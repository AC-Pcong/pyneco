package phase1

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

// TestTCPEchoServer 测试TCP Echo服务器
func TestTCPEchoServer(t *testing.T) {
	server := NewTCPEchoServer("0") // 使用随机端口

	// 异步启动服务器
	err := server.StartAsync()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 获取实际监听的端口
	port := server.listener.Addr().(*net.TCPAddr).Port

	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 发送消息
	_, err = fmt.Fprintf(conn, "Hello\n")
	if err != nil {
		t.Fatalf("Failed to send: %v", err)
	}

	// 读取响应
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	// 验证响应
	if !strings.Contains(response, "Hello") {
		t.Errorf("Expected response containing 'Hello', got '%s'", response)
	}
}
