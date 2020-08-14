package tcp

import "github.com/andycai/void"

// 带有RPC和relay功能
type MsgHooker struct {
}

func (h MsgHooker) OnInboundEvent(inputEvent void.Event) (outputEvent void.Event) {
	//var handled bool
	//var err error
	//
	//inputEvent, handled, err = rpc.ResolveInboundEvent(inputEvent)
	//
	//if err != nil {
	//	//log.Errorln("rpc.ResolveInboundEvent:", err)
	//	return
	//}
	//
	//if !handled {
	//
	//	inputEvent, handled, err = relay.ResoleveInboundEvent(inputEvent)
	//
	//	if err != nil {
	//		//log.Errorln("relay.ResoleveInboundEvent:", err)
	//		return
	//	}
	//
	//	if !handled {
	//		//msglog.WriteRecvLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())
	//	}
	//}led bool
	//var err error
	//
	//inputEvent, handled, err = rpc.ResolveInboundEvent(inputEvent)
	//
	//if err != nil {
	//	//log.Errorln("rpc.ResolveInboundEvent:", err)
	//	return
	//}
	//
	//if !handled {
	//
	//	inputEvent, handled, err = relay.ResoleveInboundEvent(inputEvent)
	//
	//	if err != nil {
	//		//log.Errorln("relay.ResoleveInboundEvent:", err)
	//		return
	//	}
	//
	//	if !handled {
	//		//msglog.WriteRecvLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())
	//	}
	//}led bool
	//var err error
	//
	//inputEvent, handled, err = rpc.ResolveInboundEvent(inputEvent)
	//
	//if err != nil {
	//	//log.Errorln("rpc.ResolveInboundEvent:", err)
	//	return
	//}
	//
	//if !handled {
	//
	//	inputEvent, handled, err = relay.ResoleveInboundEvent(inputEvent)
	//
	//	if err != nil {
	//		//log.Errorln("relay.ResoleveInboundEvent:", err)
	//		return
	//	}
	//
	//	if !handled {
	//		//msglog.WriteRecvLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())
	//	}
	//}

	return inputEvent
}

func (h MsgHooker) OnOutboundEvent(inputEvent void.Event) (outputEvent void.Event) {
	//handled, err := rpc.ResolveOutboundEvent(inputEvent)
	//
	//if err != nil {
	//	//log.Errorln("rpc.ResolveOutboundEvent:", err)
	//	return nil
	//}
	//
	//if !handled {
	//
	//	handled, err = relay.ResolveOutboundEvent(inputEvent)
	//
	//	if err != nil {
	//		//log.Errorln("relay.ResolveOutboundEvent:", err)
	//		return nil
	//	}
	//
	//	if !handled {
	//		//msglog.WriteSendLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())
	//	}
	//}

	return inputEvent
}
