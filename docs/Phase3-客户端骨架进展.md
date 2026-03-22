# Phase 3：最小客户端骨架进展

## 1. 当前目标

实现一个最小 Windows CLI 客户端，只包含：

- `connect`
- `disconnect`
- `status`

## 2. 已创建结构

- `go.mod`
- `cmd/client/main.go`
- `internal/config/config.go`
- `internal/wg/wg.go`
- `internal/status/status.go`
- `client.example.json`
- `output/build/game-booster-client.exe`

## 3. 当前实现内容

### 配置读取

当前最小 JSON 配置支持以下字段：

- `serverEndpoint`
- `serverVPNIP`
- `tunnelName`
- `configPath`
- `wireGuardPath`（可选）

配置加载逻辑已支持：

- 必填字段校验
- `serverVPNIP` 默认回退到 `10.99.0.1`
- 相对路径转绝对路径

### CLI 命令

当前 CLI 已具备以下命令入口：

- `connect`
- `disconnect`
- `status`

### WireGuard 调用方式

当前实现依赖 `WireGuard for Windows` 官方客户端，并支持以下查找顺序：

1. 配置中的 `wireGuardPath`
2. 系统 `PATH` 中的 `wireguard.exe`
3. `Program Files` 默认安装路径

命令调用方式：

- `wireguard.exe /installtunnelservice <configPath>`
- `wireguard.exe /uninstalltunnelservice <tunnelName>`

### 状态检测

当前状态检测逻辑为：

1. 查询 Windows 服务 `WireGuardTunnel$<tunnelName>`
2. 根据服务状态返回：
   - `未连接`
   - `连接中`
   - `已连接`
   - `连接失败`
3. 当服务处于 `RUNNING` 时，再 `ping 10.99.0.1` 做一次最小连通性确认

### 构建结果

已使用远程 Ubuntu 构建机交叉编译 Windows 可执行文件：

- `output/build/game-booster-client.exe`

## 4. 当前限制

- 还未在真实 Windows 客户端上执行 `connect / disconnect / status` 联调。
- 还未验证本机是否已安装 `WireGuard for Windows`。
- 还未确认云安全组是否对外放行 `51820/udp`。
- 当前尚未加入日志文件输出，仅有终端错误提示。

## 5. 下一步

下一步优先做：

1. 在 Windows 上导入 `client-01.conf` 和 `client-02.conf`
2. 用 `game-booster-client.exe` 执行 `status / connect / disconnect`
3. 验证客户端是否能访问 `10.99.0.1`
4. 验证 `10.99.0.11` 与 `10.99.0.12` 是否互通
5. 进入 Phase 4 的真实联机验证
