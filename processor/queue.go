package processor

import (
	"fmt"
	"sync"

	"github.com/andycai/void"
)

type eventQueue struct {
	*void.Pipe
	endSignal sync.WaitGroup
}

func (e *eventQueue) protectedCall(callback func()) {
	if callback != nil {
		callback()
	}
}

func (e *eventQueue) StartLoop() void.EventQueue {
	e.endSignal.Add(1)

	go func() {
		var writeList []interface{}

		for {
			writeList = writeList[0:0]
			exit := e.Pick(&writeList)

			for _, msg := range writeList {
				switch t := msg.(type) {
				case func():
					e.protectedCall(t)
				case nil:
					break
				default:
					fmt.Println("unexpected type %T", t)
				}
			}

			if exit {
				break
			}
		}

		e.endSignal.Done()
	}()

	return e
}

func (e *eventQueue) StopLoop() void.EventQueue {
	e.Add(nil)
	return e
}

func (e *eventQueue) Wait() {
	e.endSignal.Wait()
}

func (e *eventQueue) Post(callback func()) {
	if callback == nil {
		return
	}

	e.Add(callback)
}

func NewEventQueue() void.EventQueue {
	return &eventQueue{
		Pipe: void.NewPipe(),
	}
}

// 让EventCallback保证放在 session 的队列里，而不是并发的
func NewQueuedEventCallback(callback void.EventHandler) void.EventHandler {
	return func(e void.Event) {
		if callback != nil {
			SessionQueuedCall(e.Session(), func() {
				callback(e)
			})
		}
	}
}

// 在会话对应的Peer上的事件队列中执行callback，如果没有队列，则马上执行
func SessionQueuedCall(s void.Session, callback func()) {
	if s == nil {
		return
	}
	q := s.Peer().(interface {
		Queue() void.EventQueue
	}).Queue()

	QueuedCall(q, callback)
}

// 有队列时队列调用，无队列时直接调用
func QueuedCall(queue void.EventQueue, callback func()) {
	if queue == nil {
		callback()
	} else {
		queue.Post(callback)
	}
}
