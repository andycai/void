package session

import (
	"sync"
	"sync/atomic"

	"github.com/andycai/void"
)

type manager struct {
	sessionById  sync.Map // 使用Id关联会话
	sessionIDGen int64    // 记录已经生成的会话ID流水号
	count        int64    // 记录当前在使用的会话数量
}

func NewManager() void.SessionManager {
	return &manager{}
}

func (m *manager) SetIDBase(base int64) {
	atomic.StoreInt64(&m.sessionIDGen, base)
}

func (m *manager) Count() int {
	return int(atomic.LoadInt64(&m.count))
}

func (m *manager) Add(s void.Session) {
	id := atomic.AddInt64(&m.sessionIDGen, 1)
	atomic.AddInt64(&m.count, 1)
	s.SetID(id)

	m.sessionById.Store(id, s)
}

func (m *manager) Remove(s void.Session) {
	m.sessionById.Delete(s.ID())

	atomic.AddInt64(&m.count, -1)
}

func (m *manager) GetSession(id int64) void.Session {
	if v, ok := m.sessionById.Load(id); ok {
		return v.(void.Session)
	}

	return nil
}

func (m *manager) VisitSession(callback func(void.Session) bool) {
	m.sessionById.Range(func(key, value interface{}) bool {
		return callback(value.(void.Session))
	})
}

func (m *manager) CloseAllSession() {
	m.VisitSession(func(s void.Session) bool {
		s.Close()
		return true
	})
}

// 活跃的会话数量
func (m *manager) SessionCount() int {
	v := atomic.LoadInt64(&m.count)

	return int(v)
}
