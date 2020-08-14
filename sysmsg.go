package void

import "fmt"

// 系统自带消息

type SessionInit struct{}
type SessionAccepted struct{}
type SessionConnected struct{}
type SessionConnectError struct{}
type CloseReason int32

const (
	CloseReason_IO     CloseReason = iota // 普通IO断开
	CloseReason_Manual                    // 关闭前，调用过Session.Close
)

func (r CloseReason) String() string {
	switch r {
	case CloseReason_IO:
		return "IO"
	case CloseReason_Manual:
		return "Manual"
	}

	return "Unknown"
}

type SessionClosed struct {
	Reason CloseReason // 断开原因
}

// udp通知关闭,内部使用
type SessionCloseNotify struct{}

func (s *SessionInit) String() string         { return fmt.Sprintf("%+v", *s) }
func (s *SessionAccepted) String() string     { return fmt.Sprintf("%+v", *s) }
func (s *SessionConnected) String() string    { return fmt.Sprintf("%+v", *s) }
func (s *SessionConnectError) String() string { return fmt.Sprintf("%+v", *s) }
func (s *SessionClosed) String() string       { return fmt.Sprintf("%+v", *s) }
func (s *SessionCloseNotify) String() string  { return fmt.Sprintf("%+v", *s) }

// 使用类型断言判断是否为系统消息
func (s *SessionInit) SystemMessage()         {}
func (s *SessionAccepted) SystemMessage()     {}
func (s *SessionConnected) SystemMessage()    {}
func (s *SessionConnectError) SystemMessage() {}
func (s *SessionClosed) SystemMessage()       {}
func (s *SessionCloseNotify) SystemMessage()  {}
