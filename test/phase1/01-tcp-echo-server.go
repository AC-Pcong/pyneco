package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// 阶段1练习：TCP Echo服务器
// 目标：理解TCP基础、服务器监听、连接处理
func main() {
	// TODO: 实现以下功能
	// 1. 监听指定端口的TCP连接
	// 2. 接受客户端连接
	// 3. 读取客户端数据并回显
	// 4. 支持多客户端并发连接

	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	log.Printf("TCP Echo Server listening on :%s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("New connection from %s", remoteAddr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Printf("From %s: %s", remoteAddr, text)

		// Echo back
		_, err := fmt.Fprintf(conn, "Echo: %s\n", text)
		if err != nil {
			log.Printf("Write error: %v", err)
			return
		}

		// Check for quit command
		if strings.ToLower(strings.TrimSpace(text)) == "quit" {
			log.Printf("Client %s requested quit", remoteAddr)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Read error: %v", err)
	}

	log.Printf("Connection %s closed", remoteAddr)
}
