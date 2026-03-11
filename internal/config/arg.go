package config

import (
	"flag"
	"io"
	"os"
	"pushnotification_services/internal/utilities"
)

func GetDevelopmentFlavours() (string, error) {
	var envFile string

	// 使用本地 flag 集避免与 'go test' 标志冲突
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // 静默忽略未知标志（如 -test.*）
	mode := fs.String("m", "debug", "运行应用的模式 [debug|release]")

	// 解析 os.Args[1:] 查找 -m，遇到未知标志不报错
	_ = fs.Parse(os.Args[1:])

	switch *mode {
	case "release":
		envFile = ".env"
	case "debug":
		envFile = ".env.dev"
	default:
		utilities.Error("无效模式: %s，请使用 'debug' 或 'prod'", *mode)
	}

	return envFile, nil
}
