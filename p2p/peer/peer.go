package peer

import (
	"fmt"
	"io"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/linkchain/common/util/event"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/common/util/mclock"
	"github.com/linkchain/p2p/message"
	"github.com/linkchain/p2p/node"
	"github.com/linkchain/p2p/peer_error"
)

const (
	pingInterval        = 15 * time.Second
	BaseProtocolVersion = 5
	BaseProtocolLength  = uint64(16)
)

type capsByNameAndVersion []message.Cap

func (cs capsByNameAndVersion) Len() int      { return len(cs) }
func (cs capsByNameAndVersion) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs capsByNameAndVersion) Less(i, j int) bool {
	return cs[i].Name < cs[j].Name || (cs[i].Name == cs[j].Name && cs[i].Version < cs[j].Version)
}

// Peer represents a connected remote node.
type Peer struct {
	RW      *Conn
	running map[string]*protoRW
	log     log.Logger
	created mclock.AbsTime

	wg       sync.WaitGroup
	protoErr chan error
	closed   chan struct{}
	disc     chan peer_error.DiscReason

	// events receives message send / receive events if set
	events *event.Feed
}

// NewPeer returns a peer for testing purposes.
func NewTestPeer(id node.NodeID, name string, caps []message.Cap) *Peer {
	pipe, _ := net.Pipe()
	conn := &Conn{FD: pipe, Transport: nil, ID: id, Caps: caps, Name: name}
	peer := NewPeer(conn, nil)
	close(peer.closed) // ensures Disconnect doesn't block
	return peer
}

// ID returns the node's public key.
func (p *Peer) ID() node.NodeID {
	return p.RW.ID
}

func (p *Peer) CreateTime() mclock.AbsTime {
	return p.created
}

func (p *Peer) Events() *event.Feed {
	return p.events
}

func (p *Peer) SetEvents(events *event.Feed) {
	p.events = events
}

// Name returns the node name that the remote node advertised.
func (p *Peer) Name() string {
	return p.RW.Name
}

// Caps returns the capabilities (supported subprotocols) of the remote peer.
func (p *Peer) Caps() []message.Cap {
	// TODO: maybe return copy
	return p.RW.Caps
}

// RemoteAddr returns the remote address of the network connection.
func (p *Peer) RemoteAddr() net.Addr {
	return p.RW.FD.RemoteAddr()
}

// LocalAddr returns the local address of the network connection.
func (p *Peer) LocalAddr() net.Addr {
	return p.RW.FD.LocalAddr()
}

// Disconnect terminates the peer connection with the given reason.
// It returns immediately and does not wait until the connection is closed.
func (p *Peer) Disconnect(reason peer_error.DiscReason) {
	select {
	case p.disc <- reason:
	case <-p.closed:
	}
}

// String implements fmt.Stringer.
func (p *Peer) String() string {
	return fmt.Sprintf("Peer %x %v", p.RW.ID[:8], p.RemoteAddr())
}

// Inbound returns true if the peer is an inbound connection
func (p *Peer) Inbound() bool {
	return p.RW.Flags&InboundConn != 0
}

func NewPeer(conn *Conn, protocols []Protocol) *Peer {
	protomap := matchProtocols(protocols, conn.Caps, conn)
	p := &Peer{
		RW:       conn,
		running:  protomap,
		created:  mclock.Now(),
		disc:     make(chan peer_error.DiscReason),
		protoErr: make(chan error, len(protomap)+1), // protocols + pingLoop
		closed:   make(chan struct{}),
		log:      log.New("id", conn.ID, "conn", conn.Flags),
	}
	return p
}

func (p *Peer) Log() log.Logger {
	return p.log
}

func (p *Peer) Run() (remoteRequested bool, err error) {
	var (
		writeStart = make(chan struct{}, 1)
		writeErr   = make(chan error, 1)
		readErr    = make(chan error, 1)
		reason     peer_error.DiscReason // sent to the peer
	)
	p.wg.Add(2)
	go p.readLoop(readErr)
	go p.pingLoop()

	// Start all protocol handlers.
	writeStart <- struct{}{}
	p.startProtocols(writeStart, writeErr)

	// Wait for an error or disconnect.
loop:
	for {
		select {
		case err = <-writeErr:
			// A write finished. Allow the next write to start if
			// there was no error.
			if err != nil {
				reason = peer_error.DiscNetworkError
				break loop
			}
			writeStart <- struct{}{}
		case err = <-readErr:
			if r, ok := err.(peer_error.DiscReason); ok {
				remoteRequested = true
				reason = r
			} else {
				reason = peer_error.DiscNetworkError
			}
			break loop
		case err = <-p.protoErr:
			reason = peer_error.DiscReasonForError(err)
			break loop
		case err = <-p.disc:
			break loop
		}
	}

	close(p.closed)
	p.RW.Close(reason)
	p.wg.Wait()
	return remoteRequested, err
}

func (p *Peer) pingLoop() {
	ping := time.NewTimer(pingInterval)
	defer p.wg.Done()
	defer ping.Stop()
	for {
		select {
		case <-ping.C:
			if err := message.SendItems(p.RW, message.PingMsg, nil); err != nil {
				p.protoErr <- err
				return
			}
			ping.Reset(pingInterval)
		case <-p.closed:
			return
		}
	}
}

func (p *Peer) readLoop(errc chan<- error) {
	defer p.wg.Done()
	for {
		msg, err := p.RW.ReadMsg()
		if err != nil {
			errc <- err
			return
		}
		msg.ReceivedAt = time.Now()
		if err = p.handle(msg); err != nil {
			errc <- err
			return
		}
	}
}

func (p *Peer) handle(msg message.Msg) error {
	switch {
	case msg.Code == message.PingMsg:
		msg.Discard()
		go message.SendItems(p.RW, message.PongMsg, nil)
	case msg.Code == message.DiscMsg:
		var reason [1]peer_error.DiscReason
		// This is the last message. We don't need to discard or
		// check errors because, the connection will be closed after it.
		// rlp.Decode(msg.Payload, &reason)
		return reason[0]
	case msg.Code < BaseProtocolLength:
		// ignore other base protocol messages
		return msg.Discard()
	default:
		// it's a subprotocol message
		proto, err := p.getProto(msg.Code)
		if err != nil {
			return fmt.Errorf("msg code out of range: %v", msg.Code)
		}
		select {
		case proto.in <- msg:
			return nil
		case <-p.closed:
			return io.EOF
		}
	}
	return nil
}

// matchProtocols creates structures for matching named subprotocols.
func matchProtocols(protocols []Protocol, caps []message.Cap, RW message.MsgReadWriter) map[string]*protoRW {
	sort.Sort(capsByNameAndVersion(caps))
	offset := BaseProtocolLength
	result := make(map[string]*protoRW)

outer:
	for _, cap := range caps {
		for _, proto := range protocols {
			if proto.Name == cap.Name && proto.Version == cap.Version {
				// If an old protocol version matched, revert it
				if old := result[cap.Name]; old != nil {
					offset -= old.Length
				}
				// Assign the new match
				result[cap.Name] = &protoRW{Protocol: proto, offset: offset, in: make(chan message.Msg), w: RW}
				offset += proto.Length

				continue outer
			}
		}
	}
	return result
}

func (p *Peer) startProtocols(writeStart <-chan struct{}, writeErr chan<- error) {
	p.wg.Add(len(p.running))
	for _, proto := range p.running {
		proto := proto
		proto.closed = p.closed
		proto.wstart = writeStart
		proto.werr = writeErr
		var RW message.MsgReadWriter = proto
		if p.events != nil {
			RW = message.NewMsgEventer(RW, p.events, p.ID(), proto.Name)
		}
		p.log.Trace(fmt.Sprintf("Starting protocol %s/%d", proto.Name, proto.Version))
		go func() {
			err := proto.Run(p, RW)
			if err == nil {
				p.log.Trace(fmt.Sprintf("Protocol %s/%d returned", proto.Name, proto.Version))
				err = peer_error.ErrProtocolReturned
			} else if err != io.EOF {
				p.log.Trace(fmt.Sprintf("Protocol %s/%d failed", proto.Name, proto.Version), "err", err)
			}
			p.protoErr <- err
			p.wg.Done()
		}()
	}
}

// getProto finds the protocol responsible for handling
// the given message code.
func (p *Peer) getProto(code uint64) (*protoRW, error) {
	for _, proto := range p.running {
		if code >= proto.offset && code < proto.offset+proto.Length {
			return proto, nil
		}
	}
	return nil, peer_error.NewPeerError(peer_error.ErrInvalidMsgCode, "%d", code)
}

type protoRW struct {
	Protocol
	in     chan message.Msg // receices read messages
	closed <-chan struct{}  // receives when peer is shutting down
	wstart <-chan struct{}  // receives when write may start
	werr   chan<- error     // for write results
	offset uint64
	w      message.MsgWriter
}

func (RW *protoRW) WriteMsg(msg message.Msg) (err error) {
	if msg.Code >= RW.Length {
		return peer_error.NewPeerError(peer_error.ErrInvalidMsgCode, "not handled")
	}
	msg.Code += RW.offset
	select {
	case <-RW.wstart:
		err = RW.w.WriteMsg(msg)
		// Report write status back to Peer.run. It will initiate
		// shutdown if the error is non-nil and unblock the next write
		// otheRWise. The calling protocol code should exit for errors
		// as well but we don't want to rely on that.
		RW.werr <- err
	case <-RW.closed:
		err = fmt.Errorf("shutting down")
	}
	return err
}

func (RW *protoRW) ReadMsg() (message.Msg, error) {
	select {
	case msg := <-RW.in:
		msg.Code -= RW.offset
		return msg, nil
	case <-RW.closed:
		return message.Msg{}, io.EOF
	}
}

// PeerInfo represents a short summary of the information known about a connected
// peer. Sub-protocol independent fields are contained and initialized here, with
// protocol specifics delegated to all connected sub-protocols.
type PeerInfo struct {
	ID      string   `json:"id"`   // Unique node identifier (also the encryption key)
	Name    string   `json:"name"` // Name of the node, including client type, version, OS, custom data
	Caps    []string `json:"caps"` // Sum-protocols advertised by this particular peer
	Network struct {
		LocalAddress  string `json:"localAddress"`  // Local endpoint of the TCP data connection
		RemoteAddress string `json:"remoteAddress"` // Remote endpoint of the TCP data connection
		Inbound       bool   `json:"inbound"`
		Trusted       bool   `json:"trusted"`
		Static        bool   `json:"static"`
	} `json:"network"`
	Protocols map[string]interface{} `json:"protocols"` // Sub-protocol specific metadata fields
}

// Info gathers and returns a collection of metadata known about a peer.
func (p *Peer) Info() *PeerInfo {
	// Gather the protocol capabilities
	var caps []string
	for _, cap := range p.Caps() {
		caps = append(caps, cap.String())
	}
	// Assemble the generic peer metadata
	info := &PeerInfo{
		ID:        p.ID().String(),
		Name:      p.Name(),
		Caps:      caps,
		Protocols: make(map[string]interface{}),
	}
	info.Network.LocalAddress = p.LocalAddr().String()
	info.Network.RemoteAddress = p.RemoteAddr().String()
	info.Network.Inbound = p.RW.IS(InboundConn)
	info.Network.Trusted = p.RW.IS(TrustedConn)
	info.Network.Static = p.RW.IS(StaticDialedConn)

	// Gather all the running protocol infos
	for _, proto := range p.running {
		protoInfo := interface{}("unknown")
		if query := proto.Protocol.PeerInfo; query != nil {
			if metadata := query(p.ID()); metadata != nil {
				protoInfo = metadata
			} else {
				protoInfo = "handshake"
			}
		}
		info.Protocols[proto.Name] = protoInfo
	}
	return info
}
