# 阶段1练习项目

## 学习目标

通过实践项目巩固网络基础知识：
- TCP协议理解和应用
- UDP协议理解和应用
- Socket编程
- 并发处理

## 练习项目

### 1. TCP Echo Server (`01-tcp-echo-server.go`)

**功能**: 实现一个简单的TCP回显服务器

**运行方式**:
```bash
# 启动服务器
go run 01-tcp-echo-server.go 8080

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

### 2. UDP Chat (`02-udp-chat.go`)

**功能**: 实现一个基于UDP的简单聊天室

**运行方式**:
```bash
# 终端1 - 启动服务器
go run 02-udp-chat.go server :9000

# 终端2 - 客户端1
go run 02-udp-chat.go client localhost:9000 Alice

# 终端3 - 客户端2
go run 02-udp-chat.go client localhost:9000 Bob
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

### 3. Port Forward (`03-port-forward.go`)

**功能**: 实现TCP端口转发工具（类似nc）

**运行方式**:
```bash
# 先启动一个简单的echo服务器作为目标
go run 01-tcp-echo-server.go 3000

# 在另一个终端启动转发器
go run 03-port-forward.go 8080 localhost 3000

# 通过转发器访问
telnet localhost 8080
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
1. 添加流量统计（转发的字节数）
2. 添加连接超时机制
3. 支持UDP转发
4. 添加访问日志

---

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
