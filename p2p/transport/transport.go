package transport

import (
	"github.com/linkchain/p2p/message"
)

type Transport interface {
	// The two handshakes.
	// doEncHandshake(prv *ecdsa.PrivateKey, dialDest *discover.Node) (discover.NodeID, error)
	DoProtoHandshake(our *message.ProtoHandshake) (*message.ProtoHandshake, error)
	// The MsgReadWriter can only be used after the encryption
	// handshake has completed. The code uses conn.id to track this
	// by setting it to a non-nil value after the encryption handshake.
	message.MsgReadWriter
	// transports must provide Close because we use MsgPipe in some of
	// the tests. Closing the actual network connection doesn't do
	// anything in those tests because NsgPipe doesn't use it.
	Close(err error)
}
