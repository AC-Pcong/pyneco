# Pyneco - Linux网络游戏加速器

一个用Go语言开发的Linux平台网络游戏加速器。

## 项目目标

实现低延迟、低丢包的游戏网络优化，支持：
- TCP/UDP 流量代理
- TUN虚拟网卡透明代理
- UDP可靠性优化（KCP/QUIC）
- 智能路由选择
- 多节点加速

## 技术栈

- **语言**: Go
- **网络**: TCP/UDP, Socket, TUN/TAP
- **协议**: SOCKS5, KCP, 自定义协议
- **系统**: Linux (iptables, ip route)

## 项目结构

```
.
├── cmd/                 # 应用程序入口
│   ├── client/         # 客户端
│   └── server/         # 服务端
├── internal/           # 内部包
│   ├── proxy/          # 代理核心逻辑
│   ├── tun/            # TUN设备操作
│   ├── udp/            # UDP优化
│   ├── router/         # 路由管理
│   └── config/         # 配置管理
├── pkg/                # 公共包
├── docs/               # 文档
├── configs/            # 配置文件
├── scripts/            # 脚本
├── test/               # 测试相关
└── learning-plan.md    # 学习计划

```

## 开发阶段

- [x] 阶段0: 前置知识评估
- [ ] 阶段1: 网络基础夯实 (进行中)
- [ ] 阶段2: Linux网络编程入门
- [ ] 阶段3: 代理实现
- [ ] 阶段4: 流量捕获和路由
- [ ] 阶段5: UDP优化
- [ ] 阶段6: 完整系统整合

## 学习资源

参见 [learning-plan.md](learning-plan.md)

## License

MIT
