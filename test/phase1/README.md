# 阶段1练习项目

## 学习目标

通过实践项目巩固网络基础知识：
- TCP协议理解和应用
- UDP协议理解和应用
- Socket编程
- 并发处理

## 练习模块

本阶段的练习以Go包的形式提供，可以导入使用或运行测试。

### 1. TCP Echo Server

**文件**: `tcp_echo.go`, `tcp_echo_test.go`

**功能**: TCP回显服务器实现

**使用方式**:

```go
package main

import (
    "log"
    "time"
    "github.com/pcong/pyneco/test/phase1"
)

func main() {
    // 创建服务器
    server := phase1.NewTCPEchoServer("8080")

    // 异步启动
    err := server.StartAsync()
    if err != nil {
        log.Fatal(err)
    }

    // 运行测试
    time.Sleep(10 * time.Second)
    server.Stop()
}
```

**运行测试**:
```bash
cd test/phase1
go test -v -run TestTCPEchoServer
```

**手动测试**:
```bash
# 先运行一个简单的服务器示例
go run ../../examples/tcp_echo_example.go

# 使用telnet测试
telnet localhost 8080
```

**学习要点**:
- TCP三次握手过程
- Listen/Accept机制
- 连接的建立和关闭
- 并发处理多个连接

**验证任务**:
- [ ] 服务器能正常启动和监听
- [ ] 能同时处理多个客户端连接
- [ ] 客户端断开后服务器不会崩溃
- [ ] 使用tcpdump观察TCP握手过程：
  ```bash
  sudo tcpdump -i lo -n 'tcp port 8080'
  ```

---

### 2. UDP Chat

**文件**: `udp_chat.go`, `udp_chat_test.go`

**功能**: UDP聊天服务器和客户端实现

**使用方式**:

```go
// 服务器
server := phase1.NewUDPChatServer(":9000")
go server.Start()

// 客户端
client, _ := phase1.NewUDPChatClient("localhost:9009", "Alice")
client.Connect()
client.SendMessage("Hello")
```

**运行测试**:
```bash
cd test/phase1
go test -v -run TestUDPChat
```

**学习要点**:
- UDP无连接特性
- 数据报（Datagram）概念
- UDP广播机制
- 与TCP的区别

**验证任务**:
- [ ] 多个客户端能同时聊天
- [ ] 消息能正确广播
- [ ] 使用tcpdump观察UDP通信：
  ```bash
  sudo tcpdump -i lo -n 'udp port 9000'
  ```

**思考问题**:
1. UDP和TCP在实现上有什么不同？
2. UDP如何保证消息顺序？
3. 如果一个客户端掉线，服务器如何感知？

---

### 3. Port Forward

**文件**: `port_forward.go`, `port_forward_test.go`

**功能**: TCP端口转发器实现

**使用方式**:

```go
// 创建转发器
forwarder := phase1.NewPortForwarder("8080", "localhost", "3000")

// 异步启动
err := forwarder.StartAsync()
if err != nil {
    log.Fatal(err)
}
```

**运行测试**:
```bash
cd test/phase1
go test -v -run TestPortForwarder
```

**学习要点**:
- 双向数据转发
- 连接代理
- 并发数据传输
- 连接生命周期管理

**验证任务**:
- [ ] 能成功转发TCP连接
- [ ] 双向数据传输正常
- [ ] 连接关闭时两边都正确清理
- [ ] 能处理多个并发连接

**扩展练习**:
1. 添加流量统计（使用StatsPortForwarder）
2. 添加连接超时机制
3. 支持UDP转发
4. 添加访问日志

---

## 运行所有测试

```bash
cd test/phase1
go test -v
```

## 调试工具使用

### tcpdump
```bash
# 抓取TCP包
sudo tcpdump -i lo -n 'tcp port 8080'

# 抓取UDP包
sudo tcpdump -i lo -n 'udp port 9000'

# 保存到文件
sudo tcpdump -i lo -n -w capture.pcap 'tcp port 8080'

# 详细输出
sudo tcpdump -i lo -n -v 'tcp port 8080'
```

### netstat / ss
```bash
# 查看监听端口
ss -tuln | grep 8080

# 查看TCP连接
ss -tn | grep 8080

# 查看进程信息
ss -tulnp | grep 8080
```

### Wireshark
1. 使用tcpdump保存抓包文件
2. 用Wireshark打开分析：
   ```bash
   wireshark capture.pcap &
   ```

---

## 学习检验

完成所有练习后，你应该：
- ✅ 理解TCP的连接建立和关闭过程
- ✅ 理解UDP的无连接特性
- ✅ 能熟练使用Go的net包
- ✅ 能处理并发网络连接
- ✅ 能使用tcpdump分析网络流量

---

## 下一步

完成阶段1练习后，进入阶段2：Linux网络编程入门
- TUN/TAP虚拟设备
- iptables和路由表操作
- 网络命名空间
