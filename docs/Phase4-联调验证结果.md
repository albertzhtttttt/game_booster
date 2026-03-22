# Phase 4：本机联调与基础验证结果

## 1. 验证目标

本阶段目标是确认最小客户端不只是“代码存在”，而是已经能在真实 Windows 环境上完成：

- `status`
- `connect`
- `disconnect`
- 到服务端 VPN 地址 `10.99.0.1` 的基础连通

## 2. 本机前置条件检查结果

已确认：

- 本地客户端二进制存在：`output/build/game-booster-client.exe`
- 本地示例配置存在：`client.example.json`
- 本地客户端配置存在：`output/client-configs/client-01.conf`
- 初始时本机未安装 `WireGuard for Windows`
- 已通过 `winget` 安装 `WireGuard for Windows`

## 3. 实际联调过程

### 3.1 `status`

在未连接前执行：

- 返回结果：`未连接`

说明：

- 已修复中文 Windows 下 `sc.exe` 本地化输出导致的状态识别问题
- 现在服务不存在时可正确返回 `未连接`

### 3.2 `connect`

首次直接执行 `connect`：

- 返回结果：`Access is denied.`

结论：

- 启动隧道服务需要管理员权限

随后使用管理员权限重新执行 `connect`：

- 执行成功

### 3.3 连接后验证

连接成功后执行：

- `status` → 返回 `已连接`
- `ping 10.99.0.1` → 成功，延迟约 `12-13ms`

说明：

- 客户端已能通过 WireGuard 隧道访问服务端 VPN 地址
- 说明服务端与本机客户端的最小 VPN 通道已经打通

### 3.4 `disconnect`

使用管理员权限执行 `disconnect`：

- 执行成功

断开后再次执行：

- `status` → 返回 `未连接`

说明：

- 本机联调已经完成完整闭环

## 4. 当前结论

以下目标已确认达成：

- 客户端 `status` 可用
- 客户端 `connect` 可用
- 客户端 `disconnect` 可用
- 本机可连接到 `47.119.136.200` 提供的 WireGuard 服务端
- 本机已可访问 `10.99.0.1`

## 5. 仍未完成的验证

Phase 4 还剩以下内容：

- 第二台客户端导入 `client-02.conf`
- 验证 `10.99.0.11` 与 `10.99.0.12` 双向互通
- 使用目标游戏进行虚拟 IP 直连联机验证

## 6. 关键结论

当前项目已经完成：

- `M1`：服务端 `WireGuard` 成功运行
- `M2`：至少一台客户端可连通服务端
- `M3`：最小 CLI 客户端可执行 `connect / disconnect / status`

下一步应继续推进：

- 第二台客户端验证
- 双客户端互通验证
- 游戏联机验证
