package p2p

import (
	"errors"
	"net"
	"sync"

	"github.com/linkchain/common/util/log"
	"github.com/linkchain/p2p/netutil"
	"github.com/linkchain/p2p/node"
)

type Config struct {

	// MaxPeers is the maximum number of peers that can be
	// connected. It must be greater than zero.
	MaxPeers int

	// MaxPendingPeers is the maximum number of peers that can be pending in the
	// handshake phase, counted separately for inbound and outbound connections.
	// Zero defaults to preset values.
	MaxPendingPeers int `toml:",omitempty"`

	// DialRatio controls the ratio of inbound to dialed connections.
	// Example: a DialRatio of 2 allows 1/2 of connections to be dialed.
	// Setting DialRatio to zero defaults it to 3.
	DialRatio int `toml:",omitempty"`

	// Name sets the node name of this server.
	// Use common.MakeName to create a name that follows existing conventions.
	Name string `toml:"-"`

	// Connectivity can be restricted to certain IP networks.
	// If this option is set to a non-nil value, only hosts which match one of the
	// IP networks contained in the list are considered.
	NetRestrict *netutil.Netlist `toml:",omitempty"`

	// NodeDatabase is the path to the database containing the previously seen
	// live nodes in the network.
	NodeDatabase string `toml:",omitempty"`

	// Protocols should contain the protocols supported
	// by the server. Matching protocols are launched for
	// each peer.
	Protocols []Protocol `toml:"-"`

	// If ListenAddr is set to a non-nil address, the server
	// will listen for incoming connections.
	//
	// If the port is zero, the operating system will pick a port. The
	// ListenAddr field will be updated with the actual address when
	// the server is started.
	ListenAddr string

	// If Dialer is set to a non-nil value, the given Dialer
	// is used to dial outbound peer connections.
	Dialer NodeDialer `toml:"-"`

	// If NoDial is true, the server will not dial any peers.
	NoDial bool `toml:",omitempty"`

	// If EnableMsgEvents is set then the server will emit PeerEvents
	// whenever a message is sent to or received from a peer
	EnableMsgEvents bool

	// Logger is a custom logger to use with the p2p.Server.
	Logger log.Logger `toml:",omitempty"`
}

type Service struct {
	Config

	lock    sync.Mutex // protects running
	running bool
}

func (s *Service) Init(i interface{}) bool {
	log.Info("p2p service init...")
	return true
}

func (s *Service) Start() bool {
	log.Info("p2p service start...")
	return true
}

func (s *Service) Stop() {
	log.Info("p2p service stop...")
}

func (s *Service) Foo() {
	log.Info("FOOO")
}

// Self returns the local node's endpoint information.
func (srv *Service) Self() *node.Node {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	if !srv.running {
		return &node.Node{IP: net.ParseIP("0.0.0.0")}
	}
	return srv.makeSelf(srv.listener)
}

func (srv *Service) makeSelf(listener net.Listener) *node.Node {
	// Inbound connections disabled, use zero address.
	if listener == nil {
		return &node.Node{IP: net.ParseIP("0.0.0.0")}
	}
	// Otherwise inject the listener address too
	addr := listener.Addr().(*net.TCPAddr)
	return &node.Node{
		IP:  addr.IP,
		TCP: uint16(addr.Port),
	}
}

// SetupConn runs the handshakes and attempts to add the connection
// as a peer. It returns when the connection has been added as a peer
// or the handshakes have failed.
func (srv *Service) SetupConn(fd net.Conn, flags connFlag, dialDest *node.Node) error {
	self := srv.Self()
	if self == nil {
		return errors.New("shutdown")
	}
	c := &conn{fd: fd, transport: srv.newTransport(fd), flags: flags, cont: make(chan error)}
	err := srv.setupConn(c, flags, dialDest)
	if err != nil {
		c.close(err)
		srv.log.Trace("Setting up connection failed", "id", c.id, "err", err)
	}
	return err
}
