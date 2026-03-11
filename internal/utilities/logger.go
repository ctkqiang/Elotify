package utilities

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	APP_NAME = "PUSH_NOTIFICATION_SERVICE"
	VERSION  = "1.0.0"
)

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	VVERBOSE
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPink   = "\033[35m"
)

var (
	startTime     = time.Now()
	CurrentLevel  = INFO
	errorCallback func(string)
	statusMutex   sync.RWMutex
)

// 注册错误回调函数
func RegisterErrorCallback(cb func(string)) {
	errorCallback = cb
}

// 设置日志级别
func SetLogLevel(levelStr string) {
	level := strings.ToUpper(levelStr)
	switch level {
	case "DEBUG":
		CurrentLevel = DEBUG
	case "INFO":
		CurrentLevel = INFO
	case "WARN":
		CurrentLevel = WARN
	case "ERROR":
		CurrentLevel = ERROR
	case "VVERBOSE":
		CurrentLevel = VVERBOSE
	default:
		CurrentLevel = INFO
	}
}

// ToFloat64 安全地将任意接口值（主要是字符串或数值）转换为 float64
func ToFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		s := strings.TrimSpace(val)
		if s == "" || strings.EqualFold(s, "null") || strings.EqualFold(s, "none") {
			return 0
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}

func init() {
	SetLogLevel(os.Getenv("LOG_LEVEL"))
}

// 输出日志
func Log(level LogLevel, format string, a ...interface{}) {
	if level < CurrentLevel {
		return
	}

	levelStr := ""
	color := ""

	switch level {
	case DEBUG:
		levelStr = "调试"
		color = colorYellow
	case INFO:
		levelStr = "信息"
		color = colorBlue
	case WARN:
		levelStr = "警告"
		color = colorPink
	case ERROR:
		levelStr = "错误"
		color = colorRed
	case VVERBOSE:
		levelStr = "超详细"
		color = "" // 默认
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, a...)

	if level == ERROR && errorCallback != nil {
		errorCallback(msg)
	}

	if color != "" {
		fmt.Printf("["+APP_NAME+"] %s[%s] [%s] %s%s\n", color, timestamp, levelStr, msg, colorReset)
	} else {
		fmt.Printf("["+APP_NAME+"] [%s] [%s] %s\n", timestamp, levelStr, msg)
	}
}

// 信息日志
func Info(format string, a ...interface{}) { Log(INFO, format, a...) }

// 调试日志
func Debug(format string, a ...interface{}) { Log(DEBUG, format, a...) }

// 警告日志
func Warn(format string, a ...interface{}) { Log(WARN, format, a...) }

// 错误日志
func Error(format string, a ...interface{}) { Log(ERROR, format, a...) }

// 超详细日志
func VVerbose(format string, a ...interface{}) { Log(VVERBOSE, format, a...) }

// 获取环境变量，若不存在则返回默认值
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// 掩码处理字符串
func Mask(s string) string {
	runes := []rune(s)
	n := len(runes)

	if n <= 4 {
		return "****"
	}

	showCount := 10
	if n <= showCount {
		showCount = n / 3
	}

	return string(runes[:showCount]) + "[已掩码]"
}

// 检查环境变量文件是否存在
func CheckEnvFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		Error("严重：未找到 .env 文件，路径：%s", filePath)
		return
	}
	Warn("已确认：.env 文件存在于路径：%s", filePath)
}

// 检查当前内存使用情况
func CheckCUrrentMemory() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	toMB := func(bytes uint64) float64 {
		return float64(bytes) / 1024 / 1024
	}

	uptime := time.Since(startTime).Round(time.Second)
	numGoroutine := runtime.NumGoroutine()

	status := fmt.Sprintf(
		"--- [应用状态] 内存：堆=%-7.2fMB | 总计=%-7.2fMB | 系统=%-7.2fMB | 协程数：%-4d | 运行时间：%s（自 %s 起）",
		toMB(m.Alloc),
		toMB(m.TotalAlloc),
		toMB(m.Sys),
		numGoroutine,
		uptime.String(),
		startTime.Format("2006-01-02 15:04:05"),
	)

	Log(INFO, "%s", status)
	return status
}
