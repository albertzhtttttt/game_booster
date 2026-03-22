package main

import (
	"flag"
	"fmt"
	"os"

	"game_booster/internal/config"
	"game_booster/internal/status"
	"game_booster/internal/wg"
)

// main 提供最小 CLI 入口。
// 首期只保留 `connect`、`disconnect`、`status` 三个命令，避免为了命令行体验引入额外框架。
func main() {
	configPath := flag.String("config", "client.json", "客户端配置文件路径")
	flag.Parse()

	if flag.NArg() < 1 {
		printUsage()
		os.Exit(1)
	}

	command := flag.Arg(0)
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	runner := wg.Runner{
		TunnelName:    cfg.TunnelName,
		ConfigPath:    cfg.ConfigPath,
		WireGuardPath: cfg.WireGuardPath,
	}

	switch command {
	case "connect":
		if err := runner.Connect(); err != nil {
			fmt.Fprintf(os.Stderr, "connect 失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("连接命令已执行")
	case "disconnect":
		if err := runner.Disconnect(); err != nil {
			fmt.Fprintf(os.Stderr, "disconnect 失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("断开命令已执行")
	case "status":
		state, err := status.Detect(cfg.TunnelName, cfg.ServerVPNIP)
		if err != nil {
			fmt.Fprintf(os.Stderr, "status 查询失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(state)
	default:
		printUsage()
		os.Exit(1)
	}
}

// printUsage 输出最小使用说明，便于用户直接试运行。
func printUsage() {
	fmt.Println("用法: client -config client.json <connect|disconnect|status>")
}
