package phase1

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// UDPChatServer UDP聊天服务器
type UDPChatServer struct {
	conn    *net.UDPConn
	addr    string
	clients map[string]*net.UDPAddr
	mu      sync.RWMutex
}

// NewUDPChatServer 创建新的UDP聊天服务器
func NewUDPChatServer(addr string) *UDPChatServer {
	return &UDPChatServer{
		addr:    addr,
		clients: make(map[string]*net.UDPAddr),
	}
}

// Start 启动服务器（阻塞）
func (s *UDPChatServer) Start() error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.conn = conn
	defer conn.Close()

	log.Printf("UDP Chat Server listening on %s", s.addr)

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
		s.mu.Lock()
		addrKey := addr.String()
		s.clients[addrKey] = addr
		currentClients := make(map[string]*net.UDPAddr)
		for k, v := range s.clients {
			currentClients[k] = v
		}
		s.mu.Unlock()

		// 广播消息
		for _, clientAddr := range currentClients {
			if clientAddr.String() != addrKey {
				_, err := conn.WriteToUDP(buffer[:n], clientAddr)
				if err != nil {
					log.Printf("Failed to send to %s: %v", clientAddr, err)
				}
			}
		}
	}
}

// Stop 停止服务器
func (s *UDPChatServer) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// UDPChatClient UDP聊天客户端
type UDPChatClient struct {
	conn     *net.UDPConn
	serverAddr *net.UDPAddr
	nickname string
	stopChan chan struct{}
}

// NewUDPChatClient 创建新的UDP聊天客户端
func NewUDPChatClient(serverAddr, nickname string) (*UDPChatClient, error) {
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve address: %w", err)
	}

	return &UDPChatClient{
		serverAddr: addr,
		nickname:   nickname,
		stopChan:   make(chan struct{}),
	}, nil
}

// Connect 连接到服务器
func (c *UDPChatClient) Connect() error {
	conn, err := net.DialUDP("udp", nil, c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	c.conn = conn

	log.Printf("Connected to %s as %s", c.serverAddr, c.nickname)
	return nil
}

// SendJoinMessage 发送加入消息
func (c *UDPChatClient) SendJoinMessage() error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}
	message := fmt.Sprintf("[%s] joined the chat\n", c.nickname)
	_, err := fmt.Fprint(c.conn, message)
	return err
}

// SendLeaveMessage 发送离开消息
func (c *UDPChatClient) SendLeaveMessage() error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}
	message := fmt.Sprintf("[%s] left the chat\n", c.nickname)
	_, err := fmt.Fprint(c.conn, message)
	return err
}

// SendMessage 发送消息
func (c *UDPChatClient) SendMessage(text string) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}
	message := fmt.Sprintf("[%s] %s", c.nickname, text)
	_, err := fmt.Fprint(c.conn, message)
	return err
}

// StartReceiver 启动接收协程
func (c *UDPChatClient) StartReceiver(messageChan chan<- string) {
	go func() {
		buffer := make([]byte, 1024)
		for {
			select {
			case <-c.stopChan:
				return
			default:
				c.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
				n, err := c.conn.Read(buffer)
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					}
					return
				}
				messageChan <- string(buffer[:n])
			}
		}
	}()
}

// Stop 停止客户端
func (c *UDPChatClient) Stop() error {
	close(c.stopChan)
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetConnection 获取UDP连接（用于自定义读写）
func (c *UDPChatClient) GetConnection() *net.UDPConn {
	return c.conn
}
