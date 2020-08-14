package processor

import "github.com/andycai/void"

type ProcessorBinder func(bundle void.Processor, userCallback void.EventHandler, args ...interface{})

const (
	PPRO_TCP_LTV = "tcp.ltv"
	PPRO_UDP_LTV = "udp.ltv"
	PPRO_WS_LTV  = "ws.ltv"
	PPRO_HTTP    = "http"
)

var (
	procByName = map[string]ProcessorBinder{}
)

func Register(name string, f ProcessorBinder) {
	if _, ok := procByName[name]; ok {
		panic("duplicate peer type: " + name)
	}

	procByName[name] = f
}

func BindProcessorHandler(peer void.Peer, name string, userCallback void.EventHandler, args ...interface{}) {
	if binder, ok := procByName[name]; ok {
		p := peer.(void.Processor)
		binder(p, userCallback, args...)
	} else {
		panic("processor not found " + name)
	}
}
