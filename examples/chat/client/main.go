package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/andycai/void"
	"github.com/andycai/void/examples/chat/proto"
	"github.com/andycai/void/peer"
	"github.com/andycai/void/processor"

	_ "github.com/andycai/void/peer/tcp"
	_ "github.com/andycai/void/processor/tcp"
)

func ReadConsole(callback func(string)) {

	for {

		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}

		// 去掉读入内容的空白符
		text = strings.TrimSpace(text)

		callback(text)

	}

}

func main() {
	// 创建一个事件处理队列，整个客户端只有这一个队列处理事件，客户端属于单线程模型
	queue := processor.NewEventQueue()

	// 创建一个tcp的连接器，名称为client，连接地址为127.0.0.1:8801，将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	p := peer.NewPeer("tcp.Connector", "client", "127.0.0.1:12001", queue)

	// 设定封包收发处理的模式为tcp的ltv(Length-Type-Value), Length为封包大小，Type为消息ID，Value为消息内容
	// 并使用switch处理收到的消息
	processor.BindProcessorHandler(p, "tcp.ltv", func(ev void.Event) {
		switch msg := ev.Message().(type) {
		case *void.SessionConnected:
			fmt.Println("client connected")
		case *void.SessionConnectError:
			fmt.Println("client connected error")
		case *void.SessionClosed:
			fmt.Println("client error")
		case *proto.ChatACK:
			fmt.Printf("sid%d say: %s\n", msg.Id, msg.Content)
		}
	})

	// 开始发起到服务器的连接
	p.Start()

	// 事件队列开始循环
	queue.StartLoop()

	fmt.Println("Ready to chat!")

	// 阻塞的从命令行获取聊天输入
	ReadConsole(func(str string) {

		p.(interface {
			Session() void.Session
		}).Session().Send(&proto.ChatREQ{
			Content: str,
		})

	})

}
