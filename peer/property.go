package peer

import (
	"github.com/andycai/void"
)

type property struct {
	name  string
	queue void.EventQueue
	addr  string
}

func NewProperty() void.PeerProperty {
	return &property{}
}

func (p *property) Name() string {
	return p.name
}

func (p *property) Queue() void.EventQueue {
	return p.queue
}

func (p *property) Address() string {
	return p.addr
}

func (p *property) SetName(name string) {
	p.name = name
}

func (p *property) SetQueue(queue void.EventQueue) {
	p.queue = queue
}

func (p *property) SetAddress(addr string) {
	p.addr = addr
}
