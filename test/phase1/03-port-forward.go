package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

// 阶段1练习：TCP端口转发工具
// 目标：实现类似nc的端口转发功能
func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <listen_port> <target_host> <target_port>\n", os.Args[0])
		fmt.Println("Example: %s 8080 localhost 3000", os.Args[0])
		os.Exit(1)
	}

	listenPort := os.Args[1]
	targetHost := os.Args[2]
	targetPort := os.Args[3]

	// 监听本地端口
	listener, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", listenPort, err)
	}
	defer listener.Close()

	log.Printf("Port forwarder listening on :%s, forwarding to %s:%s", listenPort, targetHost, targetPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go handleConnection(conn, targetHost, targetPort)
	}
}

func handleConnection(clientConn net.Conn, targetHost, targetPort string) {
	defer clientConn.Close()

	clientAddr := clientConn.RemoteAddr().String()
	log.Printf("[%s] New connection", clientAddr)

	// 连接到目标服务器
	targetConn, err := net.Dial("tcp", targetHost+":"+targetPort)
	if err != nil {
		log.Printf("[%s] Failed to connect to target: %v", clientAddr, err)
		return
	}
	defer targetConn.Close()

	log.Printf("[%s] Connected to target %s:%s", clientAddr, targetHost, targetPort)

	// 双向转发
	var wg sync.WaitGroup
	wg.Add(2)

	// 客户端 -> 目标
	go func() {
		defer wg.Done()
		defer targetConn.Close()
		_, err := io.Copy(targetConn, clientConn)
		if err != nil {
			log.Printf("[%s] Client to target error: %v", clientAddr, err)
		}
	}()

	// 目标 -> 客户端
	go func() {
		defer wg.Done()
		defer clientConn.Close()
		_, err := io.Copy(clientConn, targetConn)
		if err != nil {
			log.Printf("[%s] Target to client error: %v", clientAddr, err)
		}
	}()

	wg.Wait()
	log.Printf("[%s] Connection closed", clientAddr)
}

// TODO: 扩展功能
// 1. 添加流量统计
// 2. 添加连接超时
// 3. 添加日志记录功能
// 4. 支持UDP转发
