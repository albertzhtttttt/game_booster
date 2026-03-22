# Phase 2：客户端配置生成结果

## 1. 阶段目标

本阶段目标是建立最小配置发放路径，并生成至少两份可导入 `WireGuard for Windows` 的客户端配置。

## 2. 配置模板结论

首期客户端配置模板固定为：

```ini
[Interface]
PrivateKey = <client-private-key>
Address = <client-vpn-ip>/24
DNS = 1.1.1.1

[Peer]
PublicKey = <server-public-key>
Endpoint = 47.119.136.200:51820
AllowedIPs = 10.99.0.0/24
PersistentKeepalive = 25
```

### 字段说明

- `Address`：客户端在 VPN 内的固定虚拟 IP
- `Endpoint`：服务端公网地址和端口
- `AllowedIPs = 10.99.0.0/24`：首期只把 VPN 子网流量走隧道，不接管全局网络
- `PersistentKeepalive = 25`：提升 NAT 场景稳定性
- `DNS = 1.1.1.1`：提供最小 DNS 配置，不影响核心联机逻辑

## 3. 已生成的客户端

### `client-01`

- 虚拟 IP：`10.99.0.11`
- 服务端地址：`47.119.136.200:51820`
- 服务端公钥：`mb/+e+mK8rkLv1C1NvWiH9ef4mETZkXKHOh+Ga3xois=`
- 服务端已写入对应 peer

### `client-02`

- 虚拟 IP：`10.99.0.12`
- 服务端地址：`47.119.136.200:51820`
- 服务端公钥：`mb/+e+mK8rkLv1C1NvWiH9ef4mETZkXKHOh+Ga3xois=`
- 服务端已写入对应 peer

## 4. 服务端状态

当前 `wg0` 已写入两个 peer：

- `10.99.0.11/32`
- `10.99.0.12/32`

服务端重启后状态正常，`wg show` 已可见两个 peer。

## 5. 服务器上的保存位置

以下文件已保存到服务器：

- `/etc/wireguard/clients/client-01.conf`
- `/etc/wireguard/clients/client-02.conf`
- `/etc/wireguard/clients/client-01.key`
- `/etc/wireguard/clients/client-02.key`
- `/etc/wireguard/clients/client-01.pub`
- `/etc/wireguard/clients/client-02.pub`

## 6. 下一步

下一步进入首次客户端导入与连通验证：

1. 将 `client-01.conf` 导入第一台 Windows 机器的 `WireGuard for Windows`
2. 将 `client-02.conf` 导入第二台 Windows 机器的 `WireGuard for Windows`
3. 启动两个隧道
4. 验证是否可访问 `10.99.0.1`
5. 验证 `10.99.0.11` 与 `10.99.0.12` 是否互通

## 7. 注意事项

- 这些配置文件包含私钥，不能公开泄露。
- 每个配置文件只能给对应的一台客户端使用。
- 如果云安全组未放行 `51820/udp`，客户端会握手失败。
- 当前尚未进行 Windows 侧实际导入验证，因此 `Phase 2` 中“手工导入首次连通测试”这一项仍待完成。
