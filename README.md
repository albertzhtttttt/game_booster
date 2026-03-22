# game_booster

## 项目简介

`game_booster` 是一个面向 Windows 的最小可用游戏联机客户端项目。

当前首期目标不是做完整“游戏加速器”，而是基于 `WireGuard` 打通一套**可实际使用、可验证、可交付**的联机方案，让两台或多台客户端通过虚拟 IP 进入同一联机网络，并用于游戏直连。

## 当前方案

首期固定技术路线为：

- `WireGuard` 星型组网
- Linux 服务端部署在 `47.119.136.200`
- Windows CLI 客户端负责 `connect / disconnect / status`
- 联机方式为“虚拟 IP 直连”

当前明确**不做**：

- GUI
- 账号系统
- 多节点切换
- 测速与智能选路
- 自动更新
- 广播发现与房间大厅

## 当前状态

目前已经完成：

- 服务端 `WireGuard` 已部署并运行
- 两份客户端配置已生成并写入服务端 peer
- Windows CLI 客户端已实现 `connect / disconnect / status`
- 单客户端已完成真实联调，可成功访问 `10.99.0.1`
- 最小发布目录与压缩包已整理完成
- 使用说明、运维说明、联调文档已补齐

当前仍未完成：

- 第二台客户端真实接入验证
- 两台客户端互相 `ping` 验证
- 目标游戏通过虚拟 IP 的真实联机验证
- 从零环境完整复现一次交付流程

## 快速开始

### 1. 准备发布目录

根据目标机器选择以下其中一套目录：

- `release/client-a`
- `release/client-b`

每套目录都包含：

- `game-booster-client.exe`
- `client.json`
- `client.conf`
- `docs/客户端使用说明.md`
- `docs/双端联机实测操作单.md`
- `docs/最终验收记录表.md`
- `docs/两台机器执行命令清单.md`

### 2. 安装 `WireGuard for Windows`

在目标 Windows 机器的 PowerShell 中执行：

```powershell
winget install --id WireGuard.WireGuard -e --source winget --accept-package-agreements --accept-source-agreements
```

安装后可执行：

```powershell
Test-Path "C:\Program Files\WireGuard\wireguard.exe"
```

返回 `True` 即表示安装成功。

### 3. 常用命令

在发布目录中打开 PowerShell：

```powershell
.\game-booster-client.exe -config .\client.json status
```

管理员 PowerShell 执行连接：

```powershell
.\game-booster-client.exe -config .\client.json connect
```

管理员 PowerShell 执行断开：

```powershell
.\game-booster-client.exe -config .\client.json disconnect
```

### 4. 单客户端验证

- 先执行 `status`，预期返回 `未连接`
- 以管理员权限执行 `connect`
- 再执行 `status`，预期返回 `已连接`
- 执行 `ping 10.99.0.1`
- 联机完成后执行 `disconnect`

## 仓库结构

```text
cmd/                     Go CLI 入口
internal/                配置、WireGuard 调用、状态检测
output/                  构建产物与客户端配置
release/                 可直接分发的最小发布目录与压缩包
docs/                    需求、方案、联调、运维与交付文档
client-a.json            客户端 A 运行配置
client-b.json            客户端 B 运行配置
client.example.json      示例配置
README.md                项目总览
```

## 关键文件

- `cmd/client/main.go`：CLI 入口
- `internal/config/config.go`：读取 JSON 配置
- `internal/wg/wg.go`：调用 `wireguard.exe` 建立或关闭隧道
- `internal/status/status.go`：检测连接状态
- `client-a.json`：客户端 A 配置
- `client-b.json`：客户端 B 配置
- `output/build/game-booster-client.exe`：已生成的 Windows 客户端

## 文档索引

- `docs/需求文档.md`
- `docs/技术方案.md`
- `docs/TODO.md`
- `docs/Phase1-服务端部署结果.md`
- `docs/Phase2-客户端配置结果.md`
- `docs/Phase3-客户端骨架进展.md`
- `docs/Phase4-联调验证结果.md`
- `docs/客户端使用说明.md`
- `docs/服务端运维说明.md`
- `docs/双端联机实测操作单.md`
- `docs/最终验收记录表.md`
- `docs/两台机器执行命令清单.md`
- `docs/最小发布结构与交付清单.md`
- `docs/发布目录说明.md`
- `docs/常见问题排查.md`
- `docs/后续扩展候选项.md`

## 配置说明

当前仓库已提供三份 JSON 配置：

- `client-a.json`：对应 `output/client-configs/client-01.conf`
- `client-b.json`：对应 `output/client-configs/client-02.conf`
- `client.example.json`：示例模板

默认服务端配置为：

- 服务端公网地址：`47.119.136.200:51820`
- 服务端 VPN 地址：`10.99.0.1`

## 构建说明

项目使用 `Go 1.22`。

如本机已安装 Go，可在仓库根目录执行：

```bash
go build -o output/build/game-booster-client.exe ./cmd/client
```

如需交叉编译 Windows 可执行文件，可执行：

```bash
GOOS=windows GOARCH=amd64 go build -o output/build/game-booster-client.exe ./cmd/client
```

## 当前交付结论

当前仓库已经达到“首期最小可交付”基础：

- 服务端可用
- 客户端可连接
- CLI 主流程可运行
- 发布目录与文档基本齐备

如果要宣称“联机功能全部完成”，仍需要补齐双客户端互通验证与目标游戏实测。