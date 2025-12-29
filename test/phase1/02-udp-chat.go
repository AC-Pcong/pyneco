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

// 阶段1练习：UDP聊天程序
// 目标：理解UDP无连接协议、数据报读写
func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage:\n")
		fmt.Printf("  Server: %s server :port\n", os.Args[0])
		fmt.Printf("  Client: %s client server_ip:port nickname\n", os.Args[0])
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "server":
		runServer(os.Args[2])
	case "client":
		if len(os.Args) < 4 {
			fmt.Println("Client requires nickname")
			os.Exit(1)
		}
		runClient(os.Args[2], os.Args[3])
	default:
		fmt.Println("Invalid mode. Use 'server' or 'client'")
	}
}

func runServer(addr string) {
	// TODO: 实现UDP服务器
	// 1. 监听UDP端口
	// 2. 接收客户端消息
	// 3. 广播消息给所有已连接的客户端
	// 4. 维护客户端列表

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Failed to resolve address: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer conn.Close()

	log.Printf("UDP Chat Server listening on %s", addr)

	// 维护客户端列表
	clients := make(map[string]*net.UDPAddr)

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}

		message := string(buffer[:n])
		log.Printf("Received from %s: %s", addr, message)

		// 注册客户端
		addrKey := addr.String()
		clients[addrKey] = addr

		// 广播消息
		for _, clientAddr := range clients {
			if clientAddr.String() != addrKey {
				_, err := conn.WriteToUDP(buffer[:n], clientAddr)
				if err != nil {
					log.Printf("Failed to send to %s: %v", clientAddr, err)
				}
			}
		}
	}
}

func runClient(serverAddr, nickname string) {
	// TODO: 实现UDP客户端
	// 1. 连接到服务器
	// 2. 发送消息
	// 3. 接收并显示消息
	// 4. 支持并发读写

	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to resolve address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to %s as %s", serverAddr, nickname)

	// 发送加入消息
	fmt.Fprintf(conn, "[%s] joined the chat\n", nickname)

	// 启动接收协程
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				log.Printf("Read error: %v", err)
				return
			}
			fmt.Printf("\r%s\n> ", string(buffer[:n]))
		}
	}()

	// 读取用户输入
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		if strings.ToLower(text) == "quit" {
			fmt.Fprintf(conn, "[%s] left the chat\n", nickname)
			break
		}

		// 发送消息
		message := fmt.Sprintf("[%s] %s", nickname, text)
		_, err := fmt.Fprint(conn, message)
		if err != nil {
			log.Printf("Send error: %v", err)
			break
		}

		fmt.Print("> ")
	}

	time.Sleep(100 * time.Millisecond)
	log.Println("Disconnected")
}
