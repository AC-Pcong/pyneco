package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pcong/pyneco/test/phase1"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <nickname> [server_address]\n", os.Args[0])
		fmt.Printf("Example: %s Alice localhost:9000\n", os.Args[0])
		os.Exit(1)
	}

	nickname := os.Args[1]
	serverAddr := "localhost:9000"
	if len(os.Args) > 2 {
		serverAddr = os.Args[2]
	}

	// 创建客户端
	client, err := phase1.NewUDPChatClient(serverAddr, nickname)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 连接到服务器
	err = client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Stop()

	log.Printf("Connected to %s as '%s'", serverAddr, nickname)
	log.Printf("Type messages and press Enter to send")
	log.Printf("Type 'quit' to exit")

	// 发送加入消息
	client.SendJoinMessage()

	// 启动接收协程
	messageChan := make(chan string, 100)
	client.StartReceiver(messageChan)

	go func() {
		for msg := range messageChan {
			fmt.Printf("\n[Rx] %s\n> ", msg)
		}
	}()

	// 读取用户输入并发送
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()

		if strings.ToLower(strings.TrimSpace(text)) == "quit" {
			client.SendLeaveMessage()
			log.Println("Leaving chat...")
			break
		}

		// 发送消息
		err = client.SendMessage(text)
		if err != nil {
			log.Printf("Failed to send: %v", err)
		} else {
			log.Printf("[Tx] Sent: %s", text)
		}

		time.Sleep(100 * time.Millisecond)
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Input error: %v", err)
	}

	log.Println("Client exiting")
}
