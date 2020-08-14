package void

// 此队列处理入站消息（revc）

type EventQueue interface {
	StartLoop() EventQueue
	StopLoop() EventQueue
	Wait()
	Post(callback func())
}
