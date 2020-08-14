package tcp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/andycai/void/processor"

	"github.com/andycai/void"
	"github.com/andycai/void/peer"
	"github.com/andycai/void/session"
)

const reportConnectFailedLimitTimes = 3

type tcpConnector struct {
	void.PeerProperty
	void.SessionManager
	void.Processor
	void.TCPSocketOption
	peer.Runnable

	defaultSession *tcpSession
	tryTimes       int // 尝试连接次数
	sesEndSignal   sync.WaitGroup
	reconDur       time.Duration
}

func (c *tcpConnector) Start() void.Peer {
	c.WaitStopFinished()
	if c.IsRunning() {
		return c
	}

	go c.connect(c.Address())

	return c
}

func (c *tcpConnector) Stop() {
	if !c.IsRunning() {
		return
	}

	if c.IsStopping() {
		return
	}

	c.StartStopping()

	c.defaultSession.Close()

	c.WaitStopFinished()
}

func (c *tcpConnector) connect(addr string) {
	c.SetRunning(true)
	for {
		c.tryTimes++

		conn, err := net.Dial("tcp", addr)

		c.defaultSession.setConn(conn)

		if err != nil {
			if c.tryTimes <= reportConnectFailedLimitTimes {
				fmt.Println("tcp.connect failed, err: " + err.Error())
				if c.tryTimes == reportConnectFailedLimitTimes {
					fmt.Println("continue reconnecting, but mute log, err: " + err.Error())
				}
			}

			// 没重连就退出
			if c.ReconnectDuration() == 0 || c.IsStopping() {
				c.ProcessEvent(processor.NewRecvEvent(c.defaultSession, &void.SessionConnectError{}))
				break
			}

			time.Sleep(c.ReconnectDuration())

			continue
		}

		c.sesEndSignal.Add(1)
		c.ApplySocketOption(conn)
		c.defaultSession.Start()
		c.tryTimes = 0
		c.ProcessEvent(processor.NewRecvEvent(c.defaultSession, &void.SessionConnected{}))
		c.sesEndSignal.Wait()

		c.defaultSession.setConn(nil)
		if c.IsStopping() || c.ReconnectDuration() == 0 {
			break
		}

		time.Sleep(c.ReconnectDuration())

		continue
	}

	c.SetRunning(false)
	c.EndStopping()
}

func (c *tcpConnector) IsReady() bool {
	return c.SessionCount() != 0
}

func (c *tcpConnector) Session() void.Session {
	return c.defaultSession
}

func (c *tcpConnector) SetSessionManager(raw interface{}) {
	c.SessionManager = raw.(void.SessionManager)
}

func (c *tcpConnector) ReconnectDuration() time.Duration {
	return c.reconDur
}

func (c *tcpConnector) SetReconnectDuration(d time.Duration) {
	c.reconDur = d
}

func (c *tcpConnector) Port() int {
	conn := c.defaultSession.Conn()
	if conn == nil {
		return 0
	}

	return conn.LocalAddr().(*net.TCPAddr).Port
}

func (c *tcpConnector) TypeName() string {
	return peer.PEER_TYPE_TCP_CONNECTOR
}

func init() {
	peer.Register(peer.PEER_TYPE_TCP_CONNECTOR, func() void.Peer {
		p := &tcpConnector{
			PeerProperty:    peer.NewProperty(),
			SessionManager:  session.NewManager(),
			Processor:       processor.NewProcessor(),
			TCPSocketOption: NewTCPSocketOption(),
			Runnable:        peer.NewRunnable(),
		}

		p.defaultSession = newSession(nil, p, func() {
			p.sesEndSignal.Done()
		})

		p.TCPSocketOption.Init()

		return p
	})
}
