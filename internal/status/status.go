package status

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

const (
	// StateDisconnected 表示隧道不存在或已停止。
	StateDisconnected = "未连接"
	// StateConnecting 表示系统正在启动或停止隧道，暂时不能判定最终状态。
	StateConnecting = "连接中"
	// StateConnected 表示隧道服务运行中，且已能访问服务端 VPN 地址。
	StateConnected = "已连接"
	// StateFailed 表示隧道服务存在，但未能通过最小连通性检查。
	StateFailed = "连接失败"
)

// Detect 使用最小规则探测连接状态。
// 当前实现遵循“先看隧道服务，再做一次轻量网络确认”的顺序：
// 1. 通过 `sc.exe query` 判断 WireGuard 隧道服务状态；
// 2. 如果服务处于 RUNNING，则再 `ping` 一次服务端 VPN 地址；
// 3. 根据结果返回“未连接 / 连接中 / 已连接 / 连接失败”。
// 这里额外兼容中文 Windows：服务不存在时优先根据退出码 `1060` 判定，而不是依赖英文输出文本。
func Detect(tunnelName string, serverVPNIP string) (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("status 仅支持在 Windows 客户端运行")
	}
	if tunnelName == "" {
		return "", fmt.Errorf("缺少隧道名称")
	}
	if serverVPNIP == "" {
		return "", fmt.Errorf("缺少服务端 VPN 地址")
	}

	serviceName := fmt.Sprintf("WireGuardTunnel$%s", tunnelName)
	cmd := exec.Command("sc.exe", "query", serviceName)
	output, err := cmd.CombinedOutput()
	serviceText := strings.ToUpper(string(output))

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1060 {
			return StateDisconnected, nil
		}
		if strings.Contains(serviceText, "FAILED 1060") || strings.Contains(serviceText, "DOES NOT EXIST") {
			return StateDisconnected, nil
		}
		return StateFailed, fmt.Errorf("查询隧道服务失败: %v: %s", err, string(output))
	}

	switch {
	case strings.Contains(serviceText, "START_PENDING"), strings.Contains(serviceText, "STOP_PENDING"):
		return StateConnecting, nil
	case strings.Contains(serviceText, "STOPPED"):
		return StateDisconnected, nil
	case strings.Contains(serviceText, "RUNNING"):
		if pingReachable(serverVPNIP) {
			return StateConnected, nil
		}
		return StateFailed, nil
	default:
		return StateFailed, nil
	}
}

// pingReachable 使用一次最小的 ICMP 探测确认服务端 VPN 地址是否可达。
// 这里只做单次短超时检查，避免 status 命令等待过久。
func pingReachable(serverVPNIP string) bool {
	cmd := exec.Command("ping", "-n", "1", "-w", "1000", serverVPNIP)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
