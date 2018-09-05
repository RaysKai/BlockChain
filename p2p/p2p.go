package p2p

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/linkchain/common/util/event"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/common/util/mclock"
	"github.com/linkchain/p2p/message"
	"github.com/linkchain/p2p/netutil"
	"github.com/linkchain/p2p/node"
	"github.com/linkchain/p2p/peer"
	"github.com/linkchain/p2p/peer_error"
	"github.com/linkchain/p2p/transport"
)

var errServerStopped = errors.New("server stopped")

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

	// BootstrapNodes are used to establish connectivity
	// with the rest of the network.
	BootstrapNodes []*node.Node

	// Static nodes are used as pre-configured connections which are always
	// maintained and re-connected on disconnects.
	StaticNodes []*node.Node

	// Trusted nodes are used as pre-configured connections which are always
	// allowed to connect, even above the peer limit.
	TrustedNodes []*node.Node

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
	Protocols []peer.Protocol `toml:"-"`

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

	lock sync.Mutex // protects running

	running bool

	ourHandshake *message.ProtoHandshake
	listener     net.Listener

	// These are for Peers, PeerCount (and nothing else).
	peerOp     chan peerOpFunc
	peerOpDone chan struct{}

	newTransport func(net.Conn) transport.Transport
	newPeerHook  func(*peer.Peer)

	quit          chan struct{}
	addstatic     chan *node.Node
	removestatic  chan *node.Node
	posthandshake chan *peer.Conn
	addpeer       chan *peer.Conn
	delpeer       chan peerDrop
	loopWG        sync.WaitGroup // loop, listenLoop
	peerFeed      event.Feed
	log           log.Logger
}

type peerOpFunc func(map[node.NodeID]*peer.Peer)

type peerDrop struct {
	*peer.Peer
	err       error
	requested bool // true if signaled by the peer
}

func (srv *Service) Init(i interface{}) bool {
	log.Info("p2p service init...")
	return true
}

func (srv *Service) Start() bool {
	log.Info("p2p service start...")
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return false
	}
	srv.running = true
	srv.log = srv.Config.Logger
	if srv.log == nil {
		srv.log = log.New()
	}
	srv.log.Info("Starting P2P networking")

	if srv.newTransport == nil {
		srv.newTransport = transport.NewPbfmsg
	}
	if srv.Dialer == nil {
		srv.Dialer = TCPDialer{&net.Dialer{Timeout: peer.DefaultDialTimeout}}
	}

	srv.quit = make(chan struct{})
	srv.addpeer = make(chan *peer.Conn)
	srv.delpeer = make(chan peerDrop)
	srv.posthandshake = make(chan *peer.Conn)
	srv.addstatic = make(chan *node.Node)
	srv.removestatic = make(chan *node.Node)
	srv.peerOp = make(chan peerOpFunc)
	srv.peerOpDone = make(chan struct{})

	dynPeers := srv.maxDialedConns()
	dialer := newDialState(srv.StaticNodes, srv.BootstrapNodes, nil, dynPeers, srv.NetRestrict)

	// handshake
	srv.ourHandshake = &message.ProtoHandshake{Version: peer.BaseProtocolVersion, Name: srv.Name}
	for _, p := range srv.Protocols {
		srv.ourHandshake.Caps = append(srv.ourHandshake.Caps, p.Cap())
	}
	// listen/dial
	if srv.ListenAddr != "" {
		if err := srv.startListening(); err != nil {
			return false
		}
	}
	if srv.NoDial && srv.ListenAddr == "" {
		srv.log.Warn("P2P server will be useless, neither dialing nor listening")
	}

	srv.loopWG.Add(1)
	go srv.run(dialer)
	srv.running = true

	return true
}

func (srv *Service) Stop() {
	log.Info("p2p service stop...")
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if !srv.running {
		return
	}
	srv.running = false
	if srv.listener != nil {
		// this unblocks listener Accept
		srv.listener.Close()
	}
	close(srv.quit)
	srv.loopWG.Wait()
}

func (srv *Service) startListening() error {
	// Launch the TCP listener.
	listener, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		return err
	}
	laddr := listener.Addr().(*net.TCPAddr)
	srv.ListenAddr = laddr.String()
	srv.listener = listener
	srv.loopWG.Add(1)
	go srv.listenLoop()
	// Map the TCP listening port if NAT is configured.
	//	if !laddr.IP.IsLoopback() && srv.NAT != nil {
	//		srv.loopWG.Add(1)
	//		go func() {
	//			nat.Map(srv.NAT, srv.quit, "tcp", laddr.Port, laddr.Port, "ethereum p2p")
	//			srv.loopWG.Done()
	//		}()
	//	}
	return nil
}

// Peers returns all connected peers.
func (srv *Service) Peers() []*peer.Peer {
	var ps []*peer.Peer
	select {
	// Note: We'd love to put this function into a variable but
	// that seems to cause a weird compiler error in some
	// environments.
	case srv.peerOp <- func(peers map[node.NodeID]*peer.Peer) {
		for _, p := range peers {
			ps = append(ps, p)
		}
	}:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return ps
}

// PeerCount returns the number of connected peers.
func (srv *Service) PeerCount() int {
	var count int
	select {
	case srv.peerOp <- func(ps map[node.NodeID]*peer.Peer) { count = len(ps) }:
		<-srv.peerOpDone
	case <-srv.quit:
	}
	return count
}

// AddPeer connects to the given node and maintains the connection until the
// server is shut down. If the connection fails for any reason, the server will
// attempt to reconnect the peer.
func (srv *Service) AddPeer(node *node.Node) {
	select {
	case srv.addstatic <- node:
	case <-srv.quit:
	}
}

// RemovePeer disconnects from the given node
func (srv *Service) RemovePeer(node *node.Node) {
	select {
	case srv.removestatic <- node:
	case <-srv.quit:
	}
}

// SubscribePeers subscribes the given channel to peer events
func (srv *Service) SubscribeEvents(ch chan *peer_error.PeerEvent) event.Subscription {
	return srv.peerFeed.Subscribe(ch)
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
func (srv *Service) SetupConn(fd net.Conn, flags peer.ConnFlag, dialDest *node.Node) error {
	self := srv.Self()
	if self == nil {
		return errors.New("shutdown")
	}
	c := peer.NewConn(fd, srv.newTransport, flags, make(chan error))
	err := srv.setupConn(c, flags, dialDest)
	if err != nil {
		c.Close(err)
		srv.log.Trace("Setting up connection failed", "id", c.ID, "err", err)
	}
	return err
}

func (srv *Service) setupConn(c *peer.Conn, flags peer.ConnFlag, dialDest *node.Node) error {
	// Prevent leftover pending conns from entering the handshake.
	srv.lock.Lock()
	running := srv.running
	srv.lock.Unlock()
	if !running {
		return errServerStopped
	}
	// Run the encryption handshake.
	var err error
	//	if c.id, err = c.doEncHandshake(srv.PrivateKey, dialDest); err != nil {
	//		srv.log.Trace("Failed RLPx handshake", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
	//		return err
	//	}
	clog := srv.log.New("id", c.ID, "addr", c.FD.RemoteAddr(), "conn", c.Flags)
	// For dialed connections, check that the remote public key matches.
	if dialDest != nil && c.ID != dialDest.ID {
		clog.Trace("Dialed identity mismatch", "want", c, dialDest.ID)
		return peer_error.DiscUnexpectedIdentity
	}
	err = srv.checkpoint(c, srv.posthandshake)
	if err != nil {
		clog.Trace("Rejected peer before protocol handshake", "err", err)
		return err
	}
	// Run the protocol handshake
	phs, err := c.DoProtoHandshake(srv.ourHandshake)
	if err != nil {
		clog.Trace("Failed proto handshake", "err", err)
		return err
	}
	if phs.ID != c.ID {
		clog.Trace("Wrong devp2p handshake identity", "err", phs.ID)
		return peer_error.DiscUnexpectedIdentity
	}
	c.Caps, c.Name = phs.Caps, phs.Name
	err = srv.checkpoint(c, srv.addpeer)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return err
	}
	// If the checks completed successfully, runPeer has now been
	// launched by run.
	clog.Trace("connection set up", "inbound", dialDest == nil)
	return nil
}

// checkpoint sends the conn to run, which performs the
// post-handshake checks for the stage (posthandshake, addpeer).
func (srv *Service) checkpoint(c *peer.Conn, stage chan<- *peer.Conn) error {
	select {
	case stage <- c:
	case <-srv.quit:
		return errServerStopped
	}
	select {
	case err := <-c.Cont:
		return err
	case <-srv.quit:
		return errServerStopped
	}
}

func (srv *Service) maxInboundConns() int {
	return srv.MaxPeers - srv.maxDialedConns()
}

func (srv *Service) maxDialedConns() int {
	if srv.NoDial {
		return 0
	}
	r := srv.DialRatio
	if r == 0 {
		r = peer.DefaultDialRatio
	}
	return srv.MaxPeers / r
}

type tempError interface {
	Temporary() bool
}

// listenLoop runs in its own goroutine and accepts
// inbound connections.
func (srv *Service) listenLoop() {
	defer srv.loopWG.Done()
	srv.log.Info("RLPx listener up", "self", srv.makeSelf(srv.listener))

	tokens := peer.DefaultMaxPendingPeers
	if srv.MaxPendingPeers > 0 {
		tokens = srv.MaxPendingPeers
	}
	slots := make(chan struct{}, tokens)
	for i := 0; i < tokens; i++ {
		slots <- struct{}{}
	}

	for {
		// Wait for a handshake slot before accepting.
		<-slots

		var (
			fd  net.Conn
			err error
		)
		for {
			fd, err = srv.listener.Accept()
			if tempErr, ok := err.(tempError); ok && tempErr.Temporary() {
				srv.log.Debug("Temporary read error", "err", err)
				continue
			} else if err != nil {
				srv.log.Debug("Read error", "err", err)
				return
			}
			break
		}

		// Reject connections that do not match NetRestrict.
		if srv.NetRestrict != nil {
			if tcp, ok := fd.RemoteAddr().(*net.TCPAddr); ok && !srv.NetRestrict.Contains(tcp.IP) {
				srv.log.Debug("Rejected conn (not whitelisted in NetRestrict)", "addr", fd.RemoteAddr())
				fd.Close()
				slots <- struct{}{}
				continue
			}
		}

		// fd = newMeteredConn(fd, true)
		srv.log.Trace("Accepted connection", "addr", fd.RemoteAddr())
		go func() {
			srv.SetupConn(fd, peer.InboundConn, nil)
			slots <- struct{}{}
		}()
	}
}

type dialer interface {
	newTasks(running int, peers map[node.NodeID]*peer.Peer, now time.Time) []task
	taskDone(task, time.Time)
	addStatic(*node.Node)
	removeStatic(*node.Node)
}

func (srv *Service) run(dialstate dialer) {
	defer srv.loopWG.Done()
	var (
		peers        = make(map[node.NodeID]*peer.Peer)
		inboundCount = 0
		trusted      = make(map[node.NodeID]bool, len(srv.TrustedNodes))
		taskdone     = make(chan task, peer.MaxActiveDialTasks)
		runningTasks []task
		queuedTasks  []task // tasks that can't run yet
	)
	// Put trusted nodes into a map to speed up checks.
	// Trusted peers are loaded on startup and cannot be
	// modified while the server is running.
	for _, n := range srv.TrustedNodes {
		trusted[n.ID] = true
	}

	// removes t from runningTasks
	delTask := func(t task) {
		for i := range runningTasks {
			if runningTasks[i] == t {
				runningTasks = append(runningTasks[:i], runningTasks[i+1:]...)
				break
			}
		}
	}
	// starts until max number of active tasks is satisfied
	startTasks := func(ts []task) (rest []task) {
		i := 0
		for ; len(runningTasks) < peer.MaxActiveDialTasks && i < len(ts); i++ {
			t := ts[i]
			srv.log.Trace("New dial task", "task", t)
			go func() { t.Do(srv); taskdone <- t }()
			runningTasks = append(runningTasks, t)
		}
		return ts[i:]
	}
	scheduleTasks := func() {
		// Start from queue first.
		queuedTasks = append(queuedTasks[:0], startTasks(queuedTasks)...)
		// Query dialer for new tasks and start as many as possible now.
		if len(runningTasks) < peer.MaxActiveDialTasks {
			nt := dialstate.newTasks(len(runningTasks)+len(queuedTasks), peers, time.Now())
			queuedTasks = append(queuedTasks, startTasks(nt)...)
		}
	}

running:
	for {
		scheduleTasks()

		select {
		case <-srv.quit:
			// The server was stopped. Run the cleanup logic.
			break running
		case n := <-srv.addstatic:
			// This channel is used by AddPeer to add to the
			// ephemeral static peer list. Add it to the dialer,
			// it will keep the node connected.
			srv.log.Debug("Adding static node", "node", n)
			dialstate.addStatic(n)
		case n := <-srv.removestatic:
			// This channel is used by RemovePeer to send a
			// disconnect request to a peer and begin the
			// stop keeping the node connected
			srv.log.Debug("Removing static node", "node", n)
			dialstate.removeStatic(n)
			if p, ok := peers[n.ID]; ok {
				p.Disconnect(peer_error.DiscRequested)
			}
		case op := <-srv.peerOp:
			// This channel is used by Peers and PeerCount.
			op(peers)
			srv.peerOpDone <- struct{}{}
		case t := <-taskdone:
			// A task got done. Tell dialstate about it so it
			// can update its state and remove it from the active
			// tasks list.
			srv.log.Trace("Dial task done", "task", t)
			dialstate.taskDone(t, time.Now())
			delTask(t)
		case c := <-srv.posthandshake:
			// A connection has passed the encryption handshake so
			// the remote identity is known (but hasn't been verified yet).
			if trusted[c.ID] {
				// Ensure that the trusted flag is set before checking against MaxPeers.
				c.Flags |= peer.TrustedConn
			}
			// TODO: track in-progress inbound node IDs (pre-Peer) to avoid dialing them.
			select {
			// case c.cont <- srv.encHandshakeChecks(peers, inboundCount, c):
			case <-srv.quit:
				break running
			}
		case c := <-srv.addpeer:
			// At this point the connection is past the protocol handshake.
			// Its capabilities are known and the remote identity is verified.
			err := srv.protoHandshakeChecks(peers, inboundCount, c)
			if err == nil {
				// The handshakes are done and it passed all checks.
				p := peer.NewPeer(c, srv.Protocols)
				// If message events are enabled, pass the peerFeed
				// to the peer
				if srv.EnableMsgEvents {
					p.SetEvents(&srv.peerFeed)
				}
				name := truncateName(c.Name)
				srv.log.Debug("Adding p2p peer", "name", name, "addr", c.FD.RemoteAddr(), "peers", len(peers)+1)
				go srv.runPeer(p)
				peers[c.ID] = p
				if p.Inbound() {
					inboundCount++
				}
			}
			// The dialer logic relies on the assumption that
			// dial tasks complete after the peer has been added or
			// discarded. Unblock the task last.
			select {
			case c.Cont <- err:
			case <-srv.quit:
				break running
			}
		case pd := <-srv.delpeer:
			// A peer disconnected.
			d := time.Duration(mclock.Now() - pd.CreateTime())
			pd.Log().Debug("Removing p2p peer", "duration", d, "peers", len(peers)-1, "req", pd.requested, "err", pd.err)
			delete(peers, pd.ID())
			if pd.Inbound() {
				inboundCount--
			}
		}
	}

	srv.log.Trace("P2P networking is spinning down")

	// Disconnect all peers.
	for _, p := range peers {
		p.Disconnect(peer_error.DiscQuitting)
	}
	// Wait for peers to shut down. Pending connections and tasks are
	// not handled here and will terminate soon-ish because srv.quit
	// is closed.
	for len(peers) > 0 {
		p := <-srv.delpeer
		p.Log().Trace("<-delpeer (spindown)", "remainingTasks", len(runningTasks))
		delete(peers, p.ID())
	}
}

func (srv *Service) protoHandshakeChecks(peers map[node.NodeID]*peer.Peer, inboundCount int, c *peer.Conn) error {
	// Drop connections with no matching protocols.
	if len(srv.Protocols) > 0 && countMatchingProtocols(srv.Protocols, c.Caps) == 0 {
		return peer_error.DiscUselessPeer
	}
	// Repeat the encryption handshake checks because the
	// peer set might have changed between the handshakes.
	// return srv.encHandshakeChecks(peers, inboundCount, c)

	return nil
}

// runPeer runs in its own goroutine for each peer.
// it waits until the Peer logic returns and removes
// the peer.
func (srv *Service) runPeer(p *peer.Peer) {
	if srv.newPeerHook != nil {
		srv.newPeerHook(p)
	}

	// broadcast peer add
	srv.peerFeed.Send(&peer_error.PeerEvent{
		Type: peer_error.PeerEventTypeAdd,
		Peer: p.ID(),
	})

	// run the protocol
	remoteRequested, err := p.Run()

	// broadcast peer drop
	srv.peerFeed.Send(&peer_error.PeerEvent{
		Type:  peer_error.PeerEventTypeDrop,
		Peer:  p.ID(),
		Error: err.Error(),
	})

	// Note: run waits for existing peers to be sent on srv.delpeer
	// before returning, so this send should not select on srv.quit.
	srv.delpeer <- peerDrop{p, err, remoteRequested}
}

func truncateName(s string) string {
	if len(s) > 20 {
		return s[:20] + "..."
	}
	return s
}

func countMatchingProtocols(protocols []peer.Protocol, caps []message.Cap) int {
	n := 0
	for _, cap := range caps {
		for _, proto := range protocols {
			if proto.Name == cap.Name && proto.Version == cap.Version {
				n++
			}
		}
	}
	return n
}
