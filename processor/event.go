package processor

import "github.com/andycai/void"

// 接受到消息
type RecvEvent struct {
	session void.Session
	message interface{}
}

func NewRecvEvent(s void.Session, m interface{}) *RecvEvent {
	return &RecvEvent{
		session: s,
		message: m,
	}
}

func (r *RecvEvent) Session() void.Session {
	return r.session
}

func (r *RecvEvent) Message() interface{} {
	return r.message
}

func (r *RecvEvent) Send(msg interface{}) {
	r.session.Send(msg)
}

// 会话开始发送数据事件
type SendEvent struct {
	session void.Session
	message interface{}
}

func NewSendEvent(s void.Session, m interface{}) *SendEvent {
	return &SendEvent{
		session: s,
		message: m,
	}
}

func (s *SendEvent) Session() void.Session {
	return s.session
}

func (s *SendEvent) Message() interface{} {
	return s.message
}
