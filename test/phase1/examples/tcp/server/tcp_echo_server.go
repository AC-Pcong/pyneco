package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pcong/pyneco/test/phase1"
)

func main() {
	// 默认端口8080
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	// 创建TCP Echo服务器
	server := phase1.NewTCPEchoServer(port)

	// 设置信号处理，优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器（阻塞）
	go func() {
		log.Printf("Starting TCP Echo Server on port %s...", port)
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
