package peer

import (
	_ "github.com/linkchain/common/util/log"
	"github.com/linkchain/p2p/message"
	"github.com/linkchain/p2p/node"
	"github.com/linkchain/p2p/transport"
	"net"
	"time"
)

type ConnFlag int

const (
	DefaultDialTimeout = 15 * time.Second

	// Connectivity defaults.
	MaxActiveDialTasks     = 16
	DefaultMaxPendingPeers = 50
	DefaultDialRatio       = 3
)

const (
	DynDialedConn ConnFlag = 1 << iota
	StaticDialedConn
	InboundConn
	TrustedConn
)

type Conn struct {
	FD net.Conn
	transport.Transport
	Flags ConnFlag
	Cont  chan error    // The run loop uses cont to signal errors to SetupConn.
	ID    node.NodeID   // valid after the encryption handshake
	Caps  []message.Cap // valid after the protocol handshake
	Name  string        // valid after the protocol handshake
}

func NewConn(fd net.Conn, transporter func(net.Conn) transport.Transport, flags ConnFlag, cont chan error) *Conn {
	return &Conn{FD: fd, Transport: transporter(fd), Flags: flags, Cont: make(chan error)}
}

func (c *Conn) String() string {
	s := c.Flags.String()
	if (c.ID != node.NodeID{}) {
		s += " " + c.ID.String()
	}
	s += " " + c.FD.RemoteAddr().String()
	return s
}

func (f ConnFlag) String() string {
	s := ""
	if f&TrustedConn != 0 {
		s += "-trusted"
	}
	if f&DynDialedConn != 0 {
		s += "-dyndial"
	}
	if f&StaticDialedConn != 0 {
		s += "-staticdial"
	}
	if f&InboundConn != 0 {
		s += "-inbound"
	}
	if s != "" {
		s = s[1:]
	}
	return s
}

func (c *Conn) IS(f ConnFlag) bool {
	return c.Flags&f != 0
}
