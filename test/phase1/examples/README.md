# Examples - 使用示例

这个目录包含如何使用 phase1 包的示例程序，用于手动测试和学习网络编程。

## TCP Echo Server 示例

### 启动服务器

```bash
# 使用默认端口 8080
go run tcp_echo_server.go

# 指定端口
go run tcp_echo_server.go 9000
```

### 启动客户端

```bash
# 终端2 - 连接到默认服务器
go run tcp_echo_client.go

# 终端3 - 再启动一个客户端
go run tcp_echo_client.go

# 终端4 - 再启动一个客户端
go run tcp_echo_client.go
```

### 观察TCP行为

在另一个终端使用 tcpdump 观察：

```bash
# 查看TCP三次握手
sudo tcpdump -i lo -n 'tcp port 8080' -S

# 查看更详细的TCP信息
sudo tcpdump -i lo -n 'tcp port 8080' -v

# 保存到文件用Wireshark分析
sudo tcpdump -i lo -n 'tcp port 8080' -w tcp_echo.pcap

# 查看连接建立和关闭
sudo tcpdump -i lo -n 'tcp[tcpflags] & (tcp-syn|tcp-fin|tcp-rst) != 0' and port 8080
```

**学习要点**：
- 观察多个客户端连接时的 TCP 三次握手（SYN, SYN-ACK, ACK）
- 观察数据传输过程中的序列号变化
- 观察连接关闭时的四次挥手（FIN, ACK）
- 观察同一端口如何处理多个连接（通过源端口区分）

---

## UDP Chat 示例

### 启动服务器

```bash
# 使用默认端口 9000
go run udp_chat_server.go

# 指定端口
go run udp_chat_server.go :8888
```

### 启动多个客户端

```bash
# 终端2 - Alice
go run udp_chat_client.go Alice

# 终端3 - Bob
go run udp_chat_client.go Bob

# 终端4 - Charlie
go run udp_chat_client.go Charlie
```

### 观察UDP行为

```bash
# 查看UDP数据包
sudo tcpdump -i lo -n 'udp port 9000'

# 保存到文件
sudo tcpdump -i lo -n 'udp port 9000' -w udp_chat.pcap

# 查看详细内容
sudo tcpdump -i lo -n 'udp port 9000' -v -X
```

**学习要点**：
- UDP没有连接建立过程（对比TCP的三次握手）
- 每个UDP数据包都是独立的
- 观察UDP如何处理丢包（如果有）
- 理解UDP的无连接特性

---

## 对比实验

### TCP vs UDP 延迟对比

```bash
# 终端1 - TCP服务器
go run tcp_echo_server.go

# 终端2 - UDP服务器
go run udp_chat_server.go

# 终端3 - 测试TCP延迟
time nc localhost 8080

# 终端4 - 测试UDP延迟（使用客户端）
go run udp_chat_client.go TestUser
```

### 使用 netstat/ss 观察连接状态

```bash
# 查看TCP连接
ss -tn | grep 8080

# 查看TCP监听
ss -tlnp | grep 8080

# 查看UDP监听
ss -ulnp | grep 9000

# 持续监控
watch -n 1 'ss -tn | grep 8080'
```

---

## 调试技巧

### 1. 使用多个终端

打开多个终端窗口，分别运行服务器和多个客户端，可以更清楚地看到每个端点的日志。

### 2. 使用 Wireshark

```bash
# 抓包并保存
sudo tcpdump -i lo -w capture.pcap

# 用Wireshark打开
wireshark capture.pcap &
```

### 3. 查看端口占用

```bash
# 查看特定端口
lsof -i :8080
ss -ltn | grep 8080
```

### 4. 压力测试

```bash
# 使用多个并发连接
for i in {1..10}; do
    go run tcp_echo_client.go &
done
```

---

## 学习检查清单

完成以下实验后，检查你的理解：

### TCP 相关
- [ ] 能识别三次握手的三个包（SYN, SYN-ACK, ACK）
- [ ] 能识别四次挥手的四个包（FIN, ACK）
- [ ] 理解序列号（Sequence Number）的作用
- [ ] 理解确认号（Acknowledgment Number）的作用
- [ ] 能看到窗口大小（Window Size）的变化
- [ ] 理解如何通过四元组（源IP、源端口、目标IP、目标端口）标识一个连接

### UDP 相关
- [ ] 理解UDP没有连接过程
- [ ] 能看到UDP数据包的边界
- [ ] 理解UDP的头部结构（比TCP简单得多）
- [ ] 知道UDP不保证可靠传输

---

## 下一步

完成这些示例学习后，你可以：
1. 修改代码添加自己的功能
2. 进入阶段2学习 TUN/TAP 设备
3. 研究 TCP 的各种状态转换（TIME_WAIT 等）
