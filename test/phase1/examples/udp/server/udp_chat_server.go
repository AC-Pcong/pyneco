package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pcong/pyneco/test/phase1"
)

func main() {
	// 默认端口9000
	addr := ":9000"
	if len(os.Args) > 1 {
		addr = os.Args[1]
		if !strings.HasPrefix(addr, ":") {
			addr = ":" + addr
		}
	}

	// 创建UDP聊天服务器
	server := phase1.NewUDPChatServer(addr)

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		log.Printf("Starting UDP Chat Server on %s...", addr)
		log.Printf("Press Ctrl+C to stop the server")
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 等待中断信号
	<-sigChan
	log.Println("\nShutting down server...")
	server.Stop()
	log.Println("Server stopped")
}
