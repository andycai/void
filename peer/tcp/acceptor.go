package tcp

import (
	"fmt"
	"net"
	"strings"

	"github.com/andycai/void/processor"

	"github.com/andycai/void"
	"github.com/andycai/void/peer"
	"github.com/andycai/void/session"
	"github.com/andycai/void/util"
)

type tcpAcceptor struct {
	void.PeerProperty
	void.SessionManager
	void.Processor
	void.TCPSocketOption
	peer.Runnable

	listener net.Listener
}

func (a *tcpAcceptor) Start() void.Peer {
	a.WaitStopFinished()

	if a.IsRunning() {
		return a
	}

	listener, err := util.DetectPort(a.Address(), func(addr *util.Addr, port int) (interface{}, error) {
		return net.Listen("tcp", addr.HostPortString(port))
	})
	if err != nil {
		//panic(fmt.Sprintf("tcp.listen error(%s) %v", a.Name(), err.Error()))
		fmt.Println("tcp.listen failed")
		a.SetRunning(false)

		return a
	}

	a.listener = listener.(net.Listener)

	go a.accept()

	return a
}

func (a *tcpAcceptor) Stop() {
	if !a.IsRunning() {
		return
	}

	if a.IsStopping() {
		return
	}

	a.StartStopping()
	a.listener.Close()
	a.CloseAllSession()

	a.WaitStopFinished()
}

func (a *tcpAcceptor) ListenAddress() string {
	pos := strings.Index(a.Address(), ":")
	if pos == -1 {
		return a.Address()
	}

	host := a.Address()[:pos]

	return util.JoinAddress(host, a.Port())
}

func (a *tcpAcceptor) accept() {
	a.SetRunning(true)

	for {
		conn, err := a.listener.Accept()

		if a.IsStopping() {
			break
		}

		if err != nil {
			continue
		}

		go a.onNewSession(conn)
	}

	a.SetRunning(false)
	a.EndStopping()
}

func (a *tcpAcceptor) onNewSession(conn net.Conn) {
	a.ApplySocketOption(conn)
	// TODO: 增加对象池管理
	s := newSession(conn, a, nil)
	s.Start()
	a.ProcessEvent(processor.NewRecvEvent(s, &void.SessionAccepted{}))
}

func (a *tcpAcceptor) Port() int {
	if a.listener == nil {
		return 0
	}

	return a.listener.Addr().(*net.TCPAddr).Port
}

func (a *tcpAcceptor) IsReady() bool {
	return a.IsRunning()
}

func (a *tcpAcceptor) TypeName() string {
	return peer.PEER_TYPE_TCP_ACCEPTOR
}

func init() {
	peer.Register(peer.PEER_TYPE_TCP_ACCEPTOR, func() void.Peer {
		p := &tcpAcceptor{
			PeerProperty:    peer.NewProperty(),
			SessionManager:  session.NewManager(),
			Processor:       processor.NewProcessor(),
			TCPSocketOption: NewTCPSocketOption(),
			Runnable:        peer.NewRunnable(),
		}

		return p
	})
}
