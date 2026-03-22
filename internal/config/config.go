package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// DefaultServerVPNIP 是首期服务端在 VPN 内使用的固定虚拟 IP。
	// 当本地配置未显式填写该字段时，客户端会回退到这个默认值，便于最小配置直接可用。
	DefaultServerVPNIP = "10.99.0.1"
)

// Config 表示最小客户端运行所需的本地配置。
// 首期只保留联机所需的最小字段：
// 1. `ServerEndpoint`：服务端公网地址，用于显示和校验；
// 2. `ServerVPNIP`：服务端在 VPN 内的虚拟 IP，用于状态探测；
// 3. `TunnelName`：Windows 服务中的隧道名称；
// 4. `ConfigPath`：实际导入 WireGuard 的 `.conf` 文件路径；
// 5. `WireGuardPath`：可选项，用于覆盖默认的 `wireguard.exe` 查找逻辑。
// 这样既保证后续 `connect / disconnect / status` 能工作，也避免引入超出首期目标的字段。
type Config struct {
	ServerEndpoint string `json:"serverEndpoint"`
	ServerVPNIP    string `json:"serverVPNIP,omitempty"`
	TunnelName     string `json:"tunnelName"`
	ConfigPath     string `json:"configPath"`
	WireGuardPath  string `json:"wireGuardPath,omitempty"`
}

// Load 从指定 JSON 文件加载客户端配置。
// 该函数除了负责解析 JSON，还会完成三件最小但必要的收尾工作：
// 1. 检查关键字段是否存在；
// 2. 为缺省的 `ServerVPNIP` 补默认值；
// 3. 将相对路径解析为“相对于配置文件所在目录”的绝对路径，避免从不同工作目录运行时找不到 `.conf`。
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if cfg.ServerEndpoint == "" {
		return Config{}, fmt.Errorf("配置缺少 serverEndpoint")
	}
	if cfg.TunnelName == "" {
		return Config{}, fmt.Errorf("配置缺少 tunnelName")
	}
	if cfg.ConfigPath == "" {
		return Config{}, fmt.Errorf("配置缺少 configPath")
	}
	if cfg.ServerVPNIP == "" {
		cfg.ServerVPNIP = DefaultServerVPNIP
	}

	configDir := filepath.Dir(path)
	if !filepath.IsAbs(cfg.ConfigPath) {
		cfg.ConfigPath = filepath.Clean(filepath.Join(configDir, cfg.ConfigPath))
	}
	if cfg.WireGuardPath != "" && !filepath.IsAbs(cfg.WireGuardPath) {
		cfg.WireGuardPath = filepath.Clean(filepath.Join(configDir, cfg.WireGuardPath))
	}

	return cfg, nil
}
