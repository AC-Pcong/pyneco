package phase1

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

// TestPortForwarder 测试端口转发器
func TestPortForwarder(t *testing.T) {
	// 先启动一个简单的echo服务器作为目标
	echoServer := NewTCPEchoServer("0")
	err := echoServer.StartAsync()
	if err != nil {
		t.Fatalf("Failed to start echo server: %v", err)
	}
	defer echoServer.Stop()

	time.Sleep(100 * time.Millisecond)

	// 获取echo服务器的端口
	targetPort := fmt.Sprintf("%d", echoServer.listener.Addr().(*net.TCPAddr).Port)

	// 启动转发器
	forwarder := NewPortForwarder("0", "127.0.0.1", targetPort)
	err = forwarder.StartAsync()
	if err != nil {
		t.Fatalf("Failed to start forwarder: %v", err)
	}
	defer forwarder.Stop()

	time.Sleep(100 * time.Millisecond)

	// 获取转发器的端口
	forwarderPort := fmt.Sprintf("%d", forwarder.listener.Addr().(*net.TCPAddr).Port)

	// 通过转发器连接到echo服务器
	conn, err := net.Dial("tcp", ":"+forwarderPort)
	if err != nil {
		t.Fatalf("Failed to connect to forwarder: %v", err)
	}
	defer conn.Close()

	// 发送测试消息
	_, err = fmt.Fprintf(conn, "Test message\n")
	if err != nil {
		t.Fatalf("Failed to send: %v", err)
	}

	// 读取响应
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	// 验证响应
	if !strings.Contains(response, "Test message") {
		t.Errorf("Expected response containing 'Test message', got '%s'", response)
	}
}
