package main

import (
	"fmt"

	"github.com/andycai/void"
	"github.com/andycai/void/examples/chat/proto"
	"github.com/andycai/void/peer"
	"github.com/andycai/void/processor"

	_ "github.com/andycai/void/peer/tcp"
	_ "github.com/andycai/void/processor/tcp"
)

func main() {
	queue := processor.NewEventQueue()

	p := peer.NewPeer("tcp.Acceptor", "server", "127.0.0.1:12001", queue)
	p.Start()

	processor.BindProcessorHandler(p, "tcp.ltv", func(e void.Event) {
		switch msg := e.Message().(type) {
		case *proto.ChatREQ:
			ack := proto.ChatACK{
				Content: msg.Content,
				Id:      e.Session().ID(),
			}
			fmt.Println("收到聊天信息：" + msg.Content)
			p.(void.SessionAccessor).VisitSession(func(session void.Session) bool {
				session.Send(&ack)

				return true
			})
		}
	})

	queue.StartLoop()
	queue.Wait()
}
