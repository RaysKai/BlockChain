package transport

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/linkchain/p2p/message"
	"github.com/linkchain/p2p/message/protobufmsg"
	"github.com/linkchain/p2p/node"
	"github.com/linkchain/p2p/peer_error"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"time"
)

const (
	maxUint24 = ^uint32(0) >> 8
	// total timeout for encryption handshake and protocol
	// handshake in both directions.
	handshakeTimeout = 5 * time.Second

	// This is the timeout for sending the disconnect reason.
	// This is shorter than the usual timeout because we don't want
	// to wait if the connection is known to be bad anyway.
	discWriteTimeout = 1 * time.Second

	// Maximum time allowed for reading a complete message.
	// This is effectively the amount of time a connection can be idle.
	frameReadTimeout = 30 * time.Second

	// Maximum amount of time allowed for writing a complete message.
	frameWriteTimeout = 20 * time.Second

	BaseProtocolMaxMsgSize = 2 * 1024
)

// errPlainMessageTooLarge is returned if a decompressed message length exceeds
// the allowed 24 bits (i.e. length >= 16MB).
var errPlainMessageTooLarge = errors.New("message length >= 16MB")

// rlpx is the transport protocol used by actual (non-test) connections.
// It wraps the frame encoder with locks and read/write deadlines.
type pbfmsg struct {
	fd net.Conn

	rmu, wmu sync.Mutex
	rw       *pbfFrameRW
}

func NewPbfmsg(fd net.Conn) Transport {
	fd.SetDeadline(time.Now().Add(handshakeTimeout))
	return &pbfmsg{fd: fd}
}

func (p *pbfmsg) ReadMsg() (message.Msg, error) {
	p.rmu.Lock()
	defer p.rmu.Unlock()
	p.fd.SetReadDeadline(time.Now().Add(frameReadTimeout))
	return p.rw.ReadMsg()
}

func (p *pbfmsg) WriteMsg(msg message.Msg) error {
	p.wmu.Lock()
	defer p.wmu.Unlock()
	p.fd.SetWriteDeadline(time.Now().Add(frameWriteTimeout))
	return p.rw.WriteMsg(msg)
}

func (p *pbfmsg) Close(err error) {
	p.wmu.Lock()
	defer p.wmu.Unlock()
	// Tell the remote end why we're disconnecting if possible.
	if p.rw != nil {
		if r, ok := err.(peer_error.DiscReason); ok && r != peer_error.DiscNetworkError {
			// rlpx tries to send DiscReason to disconnected peer
			// if the connection is net.Pipe (in-memory simulation)
			// it hangs forever, since net.Pipe does not implement
			// a write deadline. Because of this only try to send
			// the disconnect reason message if there is no error.
			if err := p.fd.SetWriteDeadline(time.Now().Add(discWriteTimeout)); err == nil {
				// TODO: fix me
				// message.SendItems(p.rw, message.DiscMsg, r)
			}
		}
	}
	p.fd.Close()
}

func (p *pbfmsg) DoProtoHandshake(our *message.ProtoHandshake) (their *message.ProtoHandshake, err error) {
	// Writing our handshake happens concurrently, we prefer
	// returning the handshake read error. If the remote side
	// disconnects us early with a valid reason, we should return it
	// as the error so it can be tracked elsewhere.
	werr := make(chan error, 1)
	go func() { // // TODO: fix me
		// werr <- message.Send(p.rw, message.HandshakeMsg, our) }()
		werr <- message.Send(p.rw, message.HandshakeMsg, nil)
	}()
	if their, err = readProtocolHandshake(p.rw, our); err != nil {
		<-werr // make sure the write terminates too
		return nil, err
	}
	if err := <-werr; err != nil {
		return nil, fmt.Errorf("write error: %v", err)
	}

	return their, nil
}

func readProtocolHandshake(rw message.MsgReader, our *message.ProtoHandshake) (*message.ProtoHandshake, error) {
	msg, err := rw.ReadMsg()
	if err != nil {
		return nil, err
	}
	if msg.Size > BaseProtocolMaxMsgSize {
		return nil, fmt.Errorf("message too big")
	}
	if msg.Code == message.DiscMsg {
		// Disconnect before protocol handshake is valid according to the
		// spec and we send it ourself if the posthanshake checks fail.
		// We can't return the reason directly, though, because it is echoed
		// back otherwise. Wrap it in a string instead.
		var reason [1]peer_error.DiscReason
		// rlp.Decode(msg.Payload, &reason)
		return nil, reason[0]
	}
	if msg.Code != message.HandshakeMsg {
		return nil, fmt.Errorf("expected handshake, got %x", msg.Code)
	}
	var hs message.ProtoHandshake
	//	if err := msg.Decode(&hs); err != nil {
	//		return nil, err
	//	}
	if (hs.ID == node.NodeID{}) {
		return nil, peer_error.DiscInvalidIdentity
	}
	return &hs, nil
}

var (
	// this is used in place of actual frame header data.
	// TODO: replace this when Msg contains the protocol type code.
	zeroHeader = []byte{0xC2, 0x80, 0x80}
	// sixteen zero bytes
	zero16 = make([]byte, 16)
)

// rlpxFrameRW implements a simplified version of RLPx framing.
// chunked messages are not supported and all headers are equal to
// zeroHeader.
//
// rlpxFrameRW is not safe for concurrent use from multiple goroutines.
type pbfFrameRW struct {
	conn io.ReadWriter
}

func newRLPXFrameRW(conn io.ReadWriter) *pbfFrameRW {
	return &pbfFrameRW{
		conn: conn,
	}
}

func (rw *pbfFrameRW) WriteMsg(msg message.Msg) error {
	content, err := ioutil.ReadAll(msg.Payload)
	if err != nil {
		return err
	}

	protobufMsg := &protobufmsg.Msg{Code: &msg.Code, Payload: content}

	data, err := proto.Marshal(protobufMsg)
	if err != nil {
		return err
	}

	headbuf := make([]byte, 32)
	dataSize := len(data)
	if uint32(dataSize) > maxUint24 {
		return errors.New("message size overflows uint24")
	}

	putInt24(uint32(dataSize), headbuf)

	if _, err := rw.conn.Write(headbuf); err != nil {
		return err
	}

	if _, err := rw.conn.Write(data); err != nil {
		return err
	}

	return nil
}

func (rw *pbfFrameRW) ReadMsg() (msg message.Msg, err error) {
	// read the header
	headbuf := make([]byte, 32)
	if _, err := io.ReadFull(rw.conn, headbuf); err != nil {
		return msg, err
	}

	dataSize := readInt24(headbuf)

	framebuf := make([]byte, dataSize)
	if _, err := io.ReadFull(rw.conn, framebuf); err != nil {
		return msg, err
	}

	protubufMsg := protobufmsg.Msg{}

	if err := proto.Unmarshal(framebuf, &protubufMsg); err != nil {
		return msg, err
	}

	msg.Code = *protubufMsg.Code
	msg.Size = uint32(len(protubufMsg.Payload))
	msg.Payload = bytes.NewReader(protubufMsg.Payload)
	msg.ReceivedAt = time.Now()

	return msg, nil
}

func readInt24(b []byte) uint32 {
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}

func putInt24(v uint32, b []byte) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}
