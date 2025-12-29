package phase1

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// PortForwarder TCP端口转发器
type PortForwarder struct {
	listenPort string
	targetHost string
	targetPort string
	listener   net.Listener
}

// NewPortForwarder 创建新的端口转发器
func NewPortForwarder(listenPort, targetHost, targetPort string) *PortForwarder {
	return &PortForwarder{
		listenPort: listenPort,
		targetHost: targetHost,
		targetPort: targetPort,
	}
}

// Start 启动转发器（阻塞）
func (p *PortForwarder) Start() error {
	listener, err := net.Listen("tcp", ":"+p.listenPort)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", p.listenPort, err)
	}
	p.listener = listener
	defer listener.Close()

	log.Printf("Port forwarder listening on :%s, forwarding to %s:%s", p.listenPort, p.targetHost, p.targetPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go p.handleConnection(conn)
	}
}

// StartAsync 异步启动转发器
func (p *PortForwarder) StartAsync() error {
	listener, err := net.Listen("tcp", ":"+p.listenPort)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", p.listenPort, err)
	}
	p.listener = listener
	log.Printf("Port forwarder listening on :%s, forwarding to %s:%s", p.listenPort, p.targetHost, p.targetPort)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Accept error: %v", err)
				return
			}
			go p.handleConnection(conn)
		}
	}()

	return nil
}

// Stop 停止转发器
func (p *PortForwarder) Stop() error {
	if p.listener != nil {
		return p.listener.Close()
	}
	return nil
}

// handleConnection 处理单个连接的转发
func (p *PortForwarder) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	clientAddr := clientConn.RemoteAddr().String()
	log.Printf("[%s] New connection", clientAddr)

	// 连接到目标服务器
	targetConn, err := net.Dial("tcp", p.targetHost+":"+p.targetPort)
	if err != nil {
		log.Printf("[%s] Failed to connect to target: %v", clientAddr, err)
		return
	}
	defer targetConn.Close()

	log.Printf("[%s] Connected to target %s:%s", clientAddr, p.targetHost, p.targetPort)

	// 双向转发
	var wg sync.WaitGroup
	wg.Add(2)

	// 客户端 -> 目标
	go func() {
		defer wg.Done()
		defer targetConn.Close()
		n, err := io.Copy(targetConn, clientConn)
		if err != nil {
			log.Printf("[%s] Client to target error: %v", clientAddr, err)
		}
		log.Printf("[%s] Client to target: %d bytes", clientAddr, n)
	}()

	// 目标 -> 客户端
	go func() {
		defer wg.Done()
		defer clientConn.Close()
		n, err := io.Copy(clientConn, targetConn)
		if err != nil {
			log.Printf("[%s] Target to client error: %v", clientAddr, err)
		}
		log.Printf("[%s] Target to client: %d bytes", clientAddr, n)
	}()

	wg.Wait()
	log.Printf("[%s] Connection closed", clientAddr)
}

// ConnectionStats 连接统计信息
type ConnectionStats struct {
	ClientAddr string
	BytesIn    int64
	BytesOut   int64
	Duration   int64 // 毫秒
}

// StatsPortForwarder 带统计功能的端口转发器
type StatsPortForwarder struct {
	*PortForwarder
	statsChan chan ConnectionStats
}

// NewStatsPortForwarder 创建带统计功能的转发器
func NewStatsPortForwarder(listenPort, targetHost, targetPort string) *StatsPortForwarder {
	return &StatsPortForwarder{
		PortForwarder: NewPortForwarder(listenPort, targetHost, targetPort),
		statsChan:     make(chan ConnectionStats, 100),
	}
}

// GetStatsChan 获取统计信息通道
func (p *StatsPortForwarder) GetStatsChan() <-chan ConnectionStats {
	return p.statsChan
}
