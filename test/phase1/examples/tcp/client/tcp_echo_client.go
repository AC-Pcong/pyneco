package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// 默认连接 localhost:8080
	serverAddr := "localhost:8080"
	if len(os.Args) > 1 {
		serverAddr = os.Args[1]
	}

	// 连接到服务器
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", serverAddr, err)
	}
	defer conn.Close()

	clientAddr := conn.LocalAddr().String()
	log.Printf("Connected to server from %s", clientAddr)
	log.Printf("Type messages and press Enter to send")
	log.Printf("Type 'quit' to exit")

	// 启动接收协程
	go func() {
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Connection closed by server")
				os.Exit(0)
			}
			fmt.Printf("\n[Rx] %s", strings.TrimSpace(line))
			fmt.Print("> ")
		}
	}()

	// 读取用户输入并发送
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()

		// 发送消息到服务器
		_, err := fmt.Fprintf(conn, "%s\n", text)
		if err != nil {
			log.Printf("Failed to send: %v", err)
			break
		}

		log.Printf("[Tx] Sent: %s", text)

		// 检查退出命令
		if strings.ToLower(strings.TrimSpace(text)) == "quit" {
			log.Println("Disconnecting...")
			break
		}

		time.Sleep(100 * time.Millisecond)
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Input error: %v", err)
	}

	log.Println("Client exiting")
}
