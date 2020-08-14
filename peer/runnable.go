package peer

import (
	"sync"
	"sync/atomic"
)

type Runnable interface {
	IsRunning() bool
	SetRunning(b bool)
	WaitStopFinished()
	IsStopping() bool
	StartStopping()
	EndStopping()
}

func NewRunnable() Runnable {
	return &runnable{
		stoppingWG: sync.WaitGroup{},
	}
}

type runnable struct {
	// 运行状态
	running int64

	stoppingWG sync.WaitGroup
	stopping   int64
}

func (r *runnable) IsRunning() bool {
	return atomic.LoadInt64(&r.running) != 0
}

func (r *runnable) SetRunning(v bool) {
	if v {
		atomic.StoreInt64(&r.running, 1)
	} else {
		atomic.StoreInt64(&r.running, 0)
	}
}

func (r *runnable) WaitStopFinished() {
	// 如果正在停止时, 等待停止完成
	r.stoppingWG.Wait()
}

func (r *runnable) IsStopping() bool {
	return atomic.LoadInt64(&r.stopping) != 0
}

func (r *runnable) StartStopping() {
	r.stoppingWG.Add(1)
	atomic.StoreInt64(&r.stopping, 1)
}

func (r *runnable) EndStopping() {
	if r.IsStopping() {
		r.stoppingWG.Done()
		atomic.StoreInt64(&r.stopping, 0)
	}
}
