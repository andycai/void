package work

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	ErrCapacity = errors.New("Thread Pool At Capacity")
)

type (
	// poolWork 被传递到工作队列里执行
	poolWork struct {
		work          PoolWorker // 被执行的工作
		resultChannel chan error // 用来通知队列操作完成
	}

	// WorkPool 实现了一个工作池，可以指定并发等级和队列容量
	WorkPool struct {
		shutdownQueueChannel chan string     // 用来关闭队列的 Channel
		shutdownWorkChannel  chan struct{}   // 用来关闭工作的 Channel
		wg                   sync.WaitGroup  // 管理关闭工作池操作
		queueChannel         chan poolWork   // 用来同步工作到队列的 Channel
		workChannel          chan PoolWorker // 用来处理工作的 Channel
		queuedWork           int32           // 已经入队的工作的数量
		activeRoutines       int32           // 正在工作的 goroutine 数量
		queueCapacity        int32           // 存储在队列中最大的容量
	}
)

// PoolWorker 必须实现此接口才能作为工作任务被执行
type PoolWorker interface {
	DoWork(workRoutine int)
}

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func NewWorkPool(numberOfRoutines int, queueCapacity int32) *WorkPool {
	p := WorkPool{
		shutdownQueueChannel: make(chan string),
		shutdownWorkChannel:  make(chan struct{}),
		queueChannel:         make(chan poolWork),
		workChannel:          make(chan PoolWorker, queueCapacity),
		queuedWork:           0,
		activeRoutines:       0,
		queueCapacity:        queueCapacity,
	}

	// 加入总的 goroutine 数量到 wait group
	p.wg.Add(numberOfRoutines)

	// 运行 goroutine 来处理工作
	for workRoutine := 0; workRoutine < numberOfRoutines; workRoutine++ {
		go p.workRoutine(workRoutine)
	}

	// 开启队列的 goroutine，用来捕捉和提供工作
	go p.queueRoutine()

	return &p
}

// Shutdown 释放资源并停止所有正在处理的工作
func (p *WorkPool) Shutdown(goRoutine string) (err error) {
	defer catchPanic(&err, goRoutine, "Shutdown")

	echo(goRoutine, "Shutdown", "Started")
	echo(goRoutine, "Shutdown", "Queue Routine")

	p.shutdownQueueChannel <- "Down"
	<-p.shutdownQueueChannel // 阻塞，保证先退出队列的 goroutine

	close(p.queueChannel)
	close(p.shutdownQueueChannel)

	echo(goRoutine, "Shutdown", "Shutting Down Work Routines")

	// 关闭所有工作 goroutine
	close(p.shutdownWorkChannel)
	// 等待所有工作 goroutine 结束
	p.wg.Wait()

	close(p.workChannel)

	echo(goRoutine, "Shutdown", "Completed")
	return err
}

// 提交工作到工作池
func (p *WorkPool) PostWork(goRoutine string, work PoolWorker) (err error) {
	defer catchPanic(&err, goRoutine, "PostWork")

	w := poolWork{work, make(chan error)}

	defer close(w.resultChannel)

	p.queueChannel <- w
	err = <-w.resultChannel // 阻塞直到有通知

	return err
}

// QueuedWork 正在执行的任务数量
func (p *WorkPool) QueuedWork() int32 {
	return atomic.AddInt32(&p.queuedWork, 0)
}

// ActiveRoutines 正在执行工作的 goroutine 数量
func (p *WorkPool) ActiveRoutines() int32 {
	return atomic.AddInt32(&p.activeRoutines, 0)
}

// CatchPanic 捕捉错误
func catchPanic(err *error, goRoutine string, functionName string) {
	if r := recover(); r != nil {
		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		echof(goRoutine, functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

func echo(goRoutine string, functionName string, message string) {
	log.Printf("%s : %s : %s\n", goRoutine, functionName, message)
}

func echof(goRoutine string, functionName string, format string, a ...interface{}) {
	echo(goRoutine, functionName, fmt.Sprintf(format, a...))
}

// 请求工作池来执行工作
func (p *WorkPool) workRoutine(workRoutine int) {
	for {
		select {
		case <-p.shutdownWorkChannel:
			echo(fmt.Sprintf("WorkRoutine %d", workRoutine), "workRoutine", "Going Down")
			p.wg.Done()
			return

		// 队列存在工作
		case poolWorker := <-p.workChannel:
			p.safelyDoWork(workRoutine, poolWorker)
			break
		}
	}
}

func (p *WorkPool) safelyDoWork(workRoutine int, poolWorker PoolWorker) {
	defer catchPanic(nil, "WorkRoutine", "SafelyDoWork")
	defer atomic.AddInt32(&p.activeRoutines, -1)

	// 更新计数
	atomic.AddInt32(&p.queuedWork, -1)
	atomic.AddInt32(&p.activeRoutines, 1)

	// 执行工作
	poolWorker.DoWork(workRoutine)
}

func (p *WorkPool) queueRoutine() {
	for {
		select {
		// 关闭队列 goroutine
		case <-p.shutdownQueueChannel:
			echo("Queue", "queueRoutine", "Going Down")
			p.shutdownQueueChannel <- "Down"
			return // 退出 goroutine

		// 提交工作准备处理
		case queueItem := <-p.queueChannel:
			// 队列已满，不允许提交
			if atomic.AddInt32(&p.queuedWork, 0) == p.queueCapacity {
				queueItem.resultChannel <- ErrCapacity
				continue
			}

			atomic.AddInt32(&p.queuedWork, 1)
			p.workChannel <- queueItem.work
			queueItem.resultChannel <- nil
			break
		}
	}
}
