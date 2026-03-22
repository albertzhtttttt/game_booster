# Phase 1：WireGuard 服务端部署结果

## 1. 部署目标

在服务器 `47.119.136.200` 上完成 `WireGuard` 最小服务端部署，使其可作为首期联机中转节点。

## 2. 实际环境

- 服务器地址：`47.119.136.200`
- SSH 用户：`root`
- 系统：`Ubuntu 24.04.3 LTS`
- 内核：`6.8.0-90-generic`

## 3. 部署结果

已完成以下内容：

- 安装 `WireGuard`
- 创建服务端配置 `/etc/wireguard/wg0.conf`
- 设置服务端虚拟地址为 `10.99.0.1/24`
- 设置监听端口为 `51820/udp`
- 启用开机自启：`wg-quick@wg0`
- 持久化开启 `net.ipv4.ip_forward = 1`
- 为 `wg0` 添加 `FORWARD` 放行规则，避免 Docker 默认转发策略影响 VPN 客户端互通

## 4. 当前关键配置

### 服务端信息

- VPN 网段：`10.99.0.0/24`
- 服务端虚拟 IP：`10.99.0.1`
- 监听端口：`51820/udp`
- 服务端公钥：`mb/+e+mK8rkLv1C1NvWiH9ef4mETZkXKHOh+Ga3xois=`

### 配置文件

当前服务端配置文件：`/etc/wireguard/wg0.conf`

核心内容如下：

```ini
[Interface]
Address = 10.99.0.1/24
ListenPort = 51820
PrivateKey = <server-private-key>
SaveConfig = false
PostUp = iptables -I FORWARD 1 -i wg0 -j ACCEPT; iptables -I FORWARD 1 -o wg0 -j ACCEPT
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -D FORWARD -o wg0 -j ACCEPT
```

## 5. 验证结果

已确认：

- `wg-quick@wg0` 服务状态为 `active`
- `wg0` 接口已创建
- `wg0` 地址为 `10.99.0.1/24`
- UDP `51820` 正在监听
- 内核 `ip_forward` 已开启
- `iptables` 的 `FORWARD` 链已插入 `wg0` 允许规则

## 6. 当前阶段结论

`Phase 1` 已完成：服务器已经具备最小可用的 `WireGuard` 服务端能力。

下一步进入 `Phase 2`：

- 设计客户端配置模板
- 生成至少两份客户端配置
- 把客户端 peer 写入服务端配置
- 手工导入 `WireGuard for Windows` 做首轮连通验证

## 7. 注意事项

- 当前系统上存在 Docker，默认 `FORWARD` 策略为 `DROP`，因此必须保留 `wg0` 的放行规则。
- 如果云平台安全组未放行 `51820/udp`，客户端仍可能无法连接，需要在云控制台同步检查。
- 当前尚未添加任何客户端 peer，故暂时只有服务端接口处于运行状态。
