package tcp

import (
	"net"
	"sync"
	"time"

	"github.com/andycai/void/processor"

	"github.com/andycai/void"
)

type tcpSession struct {
	peer      void.Peer
	conn      net.Conn
	connGuard sync.RWMutex
	sendQueue *void.Pipe // 此队列处理出站信息
	exitSync  sync.WaitGroup
	endNotify func()
	closing   bool
	id        int64
}

func newSession(conn net.Conn, p void.Peer, endNotify func()) *tcpSession {
	s := &tcpSession{
		peer:      p,
		conn:      conn,
		endNotify: endNotify,
		sendQueue: void.NewPipe(),
	}
	return s
}

func (s *tcpSession) ID() int64 {
	return s.id
}

func (s *tcpSession) SetID(id int64) {
	s.id = id
}

func (s *tcpSession) setConn(conn net.Conn) {
	s.connGuard.Lock()
	s.conn = conn
	s.connGuard.Unlock()
}

func (s *tcpSession) Conn() net.Conn {
	s.connGuard.RLock()
	defer s.connGuard.RUnlock()
	return s.conn
}

func (s *tcpSession) Peer() void.Peer {
	return s.peer
}

func (s *tcpSession) Raw() interface{} {
	return s.Conn()
}

func (s *tcpSession) Close() {
	conn := s.Conn()
	if conn != nil {
		tcpConn := conn.(*net.TCPConn)
		tcpConn.CloseRead()
		tcpConn.SetReadDeadline(time.Now())
	}
}

// Send 发送封包，加入队列等待 sendLoop 处理
func (s *tcpSession) Send(msg interface{}) {
	if msg == nil {
		return
	}

	s.sendQueue.Add(msg)
}

func (s *tcpSession) recvLoop() {
	proc := s.peer.(void.Processor)
	for s.Conn() != nil {
		msg, err := proc.ReadMessage(s)
		if err != nil {
			s.sendQueue.Add(nil) // 做标记，后面插入的数据无效

			proc.ProcessEvent(processor.NewRecvEvent(s, &void.SessionClosed{}))
			break
		}
		proc.ProcessEvent(processor.NewRecvEvent(s, msg))
	}

	s.exitSync.Done()
}

func (s *tcpSession) sendLoop() {
	proc := s.peer.(void.Processor)
	var writeList []interface{}

	for {
		writeList = writeList[0:0]
		// 没有数据时会阻塞
		exit := s.sendQueue.Pick(&writeList)

		for _, msg := range writeList {
			proc.SendMessage(processor.NewSendEvent(s, msg))
		}

		if exit {
			break
		}
	}

	conn := s.Conn()
	if conn != nil {
		conn.Close()
	}

	s.exitSync.Done()
}

func (s *tcpSession) Start() {
	mgr := s.peer.(void.SessionManager)

	s.sendQueue.Reset()

	// 接受和发送线程同时完成才算完成
	s.exitSync.Add(2)

	// 加入 session 管理器
	mgr.Add(s)

	go func() {
		s.exitSync.Wait()

		mgr.Remove(s)

		if s.endNotify != nil {
			s.endNotify()
		}
	}()

	go s.recvLoop()
	go s.sendLoop()
}
