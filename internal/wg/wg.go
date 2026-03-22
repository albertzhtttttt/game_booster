package wg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Runner 负责调用本机 WireGuard 命令。
// 这里不自己实现隧道能力，而是只做一层最小包装：
// 1. 找到 `wireguard.exe`；
// 2. 以明确参数调用官方客户端；
// 3. 将系统错误转换为更容易理解的中文提示。
type Runner struct {
	TunnelName    string
	ConfigPath    string
	WireGuardPath string
}

// Connect 启动指定隧道。
// 当前方案依赖 `WireGuard for Windows` 官方客户端，因此这里只负责调用其安装隧道服务的命令。
func (r Runner) Connect() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("connect 仅支持在 Windows 客户端运行")
	}
	if r.ConfigPath == "" {
		return fmt.Errorf("缺少配置文件路径")
	}
	if _, err := os.Stat(r.ConfigPath); err != nil {
		return fmt.Errorf("配置文件不可用: %w", err)
	}

	executablePath, err := findExecutable(r.WireGuardPath)
	if err != nil {
		return err
	}

	cmd := exec.Command(executablePath, "/installtunnelservice", r.ConfigPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return formatCommandError("启动隧道", err, output)
	}

	return nil
}

// Disconnect 关闭指定隧道。
// 隧道名称需要与 `.conf` 文件导入后的服务名一致，因此这里直接使用配置中的 `TunnelName`。
func (r Runner) Disconnect() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("disconnect 仅支持在 Windows 客户端运行")
	}
	if r.TunnelName == "" {
		return fmt.Errorf("缺少隧道名称")
	}

	executablePath, err := findExecutable(r.WireGuardPath)
	if err != nil {
		return err
	}

	cmd := exec.Command(executablePath, "/uninstalltunnelservice", r.TunnelName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return formatCommandError("关闭隧道", err, output)
	}

	return nil
}

// findExecutable 负责定位 `wireguard.exe`。
// 为了兼容常见安装方式，它按以下顺序查找：
// 1. 配置文件中显式指定的路径；
// 2. 当前 PATH 中的 `wireguard.exe`；
// 3. `Program Files` 和 `Program Files (x86)` 下的默认安装目录。
// 如果都找不到，则返回明确错误，提示用户先安装官方客户端或手动指定路径。
func findExecutable(explicitPath string) (string, error) {
	if explicitPath != "" {
		if fileInfo, err := os.Stat(explicitPath); err == nil && !fileInfo.IsDir() {
			return explicitPath, nil
		}
		return "", fmt.Errorf("配置中的 wireGuardPath 不可用: %s", explicitPath)
	}

	if executablePath, err := exec.LookPath("wireguard.exe"); err == nil {
		return executablePath, nil
	}

	candidatePaths := []string{
		filepath.Join(os.Getenv("ProgramFiles"), "WireGuard", "wireguard.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "WireGuard", "wireguard.exe"),
		`C:\Program Files\WireGuard\wireguard.exe`,
		`C:\Program Files (x86)\WireGuard\wireguard.exe`,
	}

	for _, candidatePath := range candidatePaths {
		if candidatePath == "" {
			continue
		}
		if fileInfo, err := os.Stat(candidatePath); err == nil && !fileInfo.IsDir() {
			return candidatePath, nil
		}
	}

	return "", fmt.Errorf("未找到 wireguard.exe，请先安装 WireGuard for Windows，或在配置中填写 wireGuardPath")
}

// formatCommandError 将命令执行错误转成更直接的中文提示。
// 当前重点处理管理员权限不足的场景，因为本项目在真实联调中已确认安装/卸载隧道服务需要提权。
func formatCommandError(action string, err error, output []byte) error {
	outputText := string(output)
	if strings.Contains(strings.ToLower(outputText), "access is denied") {
		return fmt.Errorf("%s失败：需要管理员权限，请使用管理员方式运行客户端或提权后重试", action)
	}
	return fmt.Errorf("%s失败: %v: %s", action, err, outputText)
}
