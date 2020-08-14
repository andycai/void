package void

import (
	"sync"
)

// 不限制大小，添加不发生阻塞，接收阻塞等待
type Pipe struct {
	list      []interface{}
	listGuard sync.Mutex
	listCond  *sync.Cond // 当修改条件或者调用Wait方法时，必须加锁
}

// Add 添加时不会发送阻塞
func (p *Pipe) Add(msg interface{}) {
	p.listGuard.Lock()
	p.list = append(p.list, msg)
	p.listGuard.Unlock()

	p.listCond.Signal() // 通知 c 最近的 Wait 继续执行，见 p.listCond.Wait()
}

func (p *Pipe) Reset() {
	p.list = p.list[0:0]
}

// Pick 从列表获取数据，并清空列表；如果没有数据，发生阻塞
func (p *Pipe) Pick(retList *[]interface{}) (exit bool) {
	p.listGuard.Lock() // cond.Wait 必须要加锁

	for len(p.list) == 0 {
		p.listCond.Wait() // 阻塞
	}

	p.listGuard.Unlock()

	p.listGuard.Lock()

	// 复制出队列
	for _, data := range p.list {
		// 有 nil 的数据，return exit标记
		if data == nil {
			exit = true
			break
		} else {
			*retList = append(*retList, data)
		}
	}

	p.Reset()
	p.listGuard.Unlock()

	return
}

func NewPipe() *Pipe {
	p := &Pipe{}
	p.listCond = sync.NewCond(&p.listGuard)

	return p
}
