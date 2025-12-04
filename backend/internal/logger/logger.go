package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Logger 全局日志实例
var Logger *logrus.Logger

// Config 日志配置
type Config struct {
	LogDir       string
	MaxAge       time.Duration
	RotationTime time.Duration
	Level        string
}

// GetDefaultConfig 获取默认日志配置
func GetDefaultConfig() *Config {
	return &Config{
		LogDir:       "./logs",
		MaxAge:       7 * 24 * time.Hour, // 保留7天
		RotationTime: 24 * time.Hour,     // 每天轮转一次
		Level:        "info",             // 默认日志级别
	}
}

// InitLogger 初始化日志
func InitLogger(config *Config) error {
	// 创建日志目录
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 创建logrus实例
	Logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 设置日志格式
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 配置日志轮转
	logPath := filepath.Join(config.LogDir, "app.log")
	rotatelogger, err := rotatelogs.New(
		logPath+".%Y%m%d",
		rotatelogs.WithMaxAge(config.MaxAge),
		rotatelogs.WithRotationTime(config.RotationTime),
	)
	if err != nil {
		return fmt.Errorf("failed to create rotatelogger: %w", err)
	}

	// 添加文件日志钩子
	hook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: rotatelogger,
		logrus.InfoLevel:  rotatelogger,
		logrus.WarnLevel:  rotatelogger,
		logrus.ErrorLevel: rotatelogger,
		logrus.FatalLevel: rotatelogger,
		logrus.PanicLevel: rotatelogger,
	}, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	Logger.AddHook(hook)

	// 同时输出到控制台
	Logger.SetOutput(os.Stdout)

	return nil
}

// Debug 调试日志
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}
