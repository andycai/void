package main

import (
	"fmt"

	"github.com/andycai/void"
	"github.com/andycai/void/peer"
	"github.com/andycai/void/processor"
)

func clientAsyncCallback() {

	// 等待服务器返回数据
	done := make(chan struct{})

	queue := processor.NewEventQueue()

	p := peer.NewPeer("tcp.Connector", "clientAsyncCallback", peerAddress, queue)

	processor.BindProcessorHandler(p, "tcp.ltv", func(ev void.Event) {

		switch msg := ev.Message().(type) {
		case *void.SessionConnected: // 已经连接上
			fmt.Println("clientAsyncCallback connected")
			ev.Session().Send(&TestEchoACK{
				Msg:   "hello",
				Value: 1234,
			})
		case *TestEchoACK: //收到服务器发送的消息

			fmt.Printf("clientAsyncCallback recv %+v\n", msg)

			// 完成操作
			done <- struct{}{}

		case *void.SessionClosed:
			fmt.Println("clientAsyncCallback closed")
		}
	})

	p.Start()

	queue.StartLoop()

	// 等待客户端收到消息
	<-done
}
