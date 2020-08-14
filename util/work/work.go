package work

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	addRoutine    = 1
	removeRoutine = 2
)

var ErrorInvalidMinRoutines = errors.New("Invalid minimum number of routines")
var ErrorInvalidStatTime = errors.New("Invalid duration for stat time")

// 必须实现 Worker 接口的任务才能被工作池管理和执行
type Worker interface {
	Work(id int)
}

// Pool 提供一个能执行工作任务的的协程池
type Pool struct {
	minRoutines int            // 池中最小的 goroutine 数量
	statTime    time.Duration  // 显示状态的时间间隔
	counter     int            // 创建 goroutine 的总数量，同时也作为创建的 goroutine ID
	tasks       chan Worker    // 没有缓冲的 work channel，用于处理任务
	control     chan int       // 没有缓冲的 work channel，用于创建和上次 goroutine
	kill        chan struct{}  // 没有缓冲的 work channel，用于杀死 goroutine
	shutdown    chan struct{}  // 关闭工作池的 channel
	wg          sync.WaitGroup // 管理关闭工作池计数
	routines    int64          // 当前 goroutine 数量
	active      int64          // 工作池中正在执行任务的 goroutine 数量
	pending     int64          // 等待执行的任务数量

	logFunc func(message string) // 日志输出函数
}

func NewWork(minRoutines int, statTime time.Duration, logFunc func(message string)) (*Pool, error) {
	if minRoutines <= 0 {
		return nil, ErrorInvalidMinRoutines
	}

	if statTime < time.Millisecond {
		return nil, ErrorInvalidStatTime
	}

	p := Pool{
		minRoutines: minRoutines,
		statTime:    statTime,
		tasks:       make(chan Worker),   // 可以带缓冲
		control:     make(chan int),      // 创建或删除 goroutine
		kill:        make(chan struct{}), // 结束 goroutine
		shutdown:    make(chan struct{}), // 关闭整个工作池
		logFunc:     logFunc,             // 日志输出回调
	}

	p.manager()
	p.Add(minRoutines)

	return &p, nil
}

func (p *Pool) Add(routines int) {
	if routines == 0 {
		return
	}

	cmd := addRoutine
	if routines < 0 {
		routines = routines * -1
		cmd = removeRoutine
	}

	for i := 0; i < routines; i++ {
		p.control <- cmd
	}
}

func (p *Pool) work(id int) {
done: // 退出循环标签
	for {
		select {
		case t := <-p.tasks:
			atomic.AddInt64(&p.active, 1)
			{
				t.Work(id)
			}
			atomic.AddInt64(&p.active, -1)

		case <-p.kill: // 退出一个 goroutine
			break done // 跳出 for 循环
		}
	}

	// 结束当前 goroutine
	atomic.AddInt64(&p.routines, -1)
	p.wg.Done() // WaitGroup 计数-1

	p.log("Worker : Shutting Down" + ", work id: " + strconv.Itoa(id))
}

func (p *Pool) Run(work Worker) {
	atomic.AddInt64(&p.pending, 1)
	{
		p.tasks <- work
	}
	atomic.AddInt64(&p.pending, -1)
}

func (p *Pool) Shutdown() {
	close(p.shutdown)
	p.wg.Wait() // 等待所有 goroutine 退出
}

func (p *Pool) manager() {
	p.wg.Add(1)

	go func() {
		p.log("Work Manager : Started")

		// 状态显示定期器
		timer := time.NewTimer(p.statTime)

		for {
			select {
			case <-p.shutdown:
				routines := int(atomic.LoadInt64(&p.routines))

				for i := 0; i < routines; i++ {
					p.kill <- struct{}{}
				}

				p.wg.Done()
				return

			case c := <-p.control:
				switch c {
				case addRoutine:
					p.log("Work Manager : Add Routine")

					p.counter++

					p.wg.Add(1)
					atomic.AddInt64(&p.routines, 1)

					go p.work(p.counter)

				case removeRoutine:
					p.log("Work Manager : Remove Routine")

					routines := int(atomic.LoadInt64(&p.routines))

					if routines <= p.minRoutines {
						p.log("Work Manager : Reached Minimum Can't Remove")
						break
					}

					p.kill <- struct{}{}
				}

			case <-timer.C:
				routines := atomic.LoadInt64(&p.routines)
				pending := atomic.LoadInt64(&p.pending)
				active := atomic.LoadInt64(&p.active)

				p.log(fmt.Sprintf("Work Manager : Stats : G[%d] P[%d] A[%d]", routines, pending, active))

				timer.Reset(p.statTime)
			}
		}
	}()
}

func (p *Pool) log(message string) {
	if p.logFunc != nil {
		p.logFunc(message)
	}
}
