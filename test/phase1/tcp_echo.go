package phase1

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// TCPEchoServer TCP Echo服务器
// 用于学习TCP协议基础、服务器监听、连接处理
type TCPEchoServer struct {
	listener net.Listener
	port     string
}

// NewTCPEchoServer 创建新的TCP Echo服务器
func NewTCPEchoServer(port string) *TCPEchoServer {
	return &TCPEchoServer{
		port: port,
	}
}

// Start 启动服务器（阻塞）
func (s *TCPEchoServer) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.listener = listener
	log.Printf("TCP Echo Server listening on :%s", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// Stop 停止服务器
func (s *TCPEchoServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// handleConnection 处理单个客户端连接
func (s *TCPEchoServer) handleConnection(conn net.Conn) {
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

// StartAsync 异步启动服务器
func (s *TCPEchoServer) StartAsync() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.listener = listener
	log.Printf("TCP Echo Server listening on :%s", s.port)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Accept error: %v", err)
				return
			}
			go s.handleConnection(conn)
		}
	}()

	return nil
}
