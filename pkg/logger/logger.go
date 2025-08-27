package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"blog-api/pkg/config"
	
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志记录器接口
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// simpleLogger 简单日志实现
type simpleLogger struct {
	level  Level
	writer io.Writer
}

// Level 日志级别
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

var globalLogger Logger

// InitLogger 初始化日志系统
func InitLogger(cfg *config.LogConfig) error {
	// 确保日志目录存在
	if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
		return err
	}

	// 配置日志轮转
	logWriter := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	// 根据环境选择输出
	var writer io.Writer
	if gin.Mode() == gin.ReleaseMode {
		writer = logWriter
	} else {
		writer = io.MultiWriter(os.Stdout, logWriter)
	}

	// 设置日志级别
	level := parseLevel(cfg.Level)
	
	globalLogger = &simpleLogger{
		level:  level,
		writer: writer,
	}

	return nil
}

// parseLevel 解析日志级别
func parseLevel(levelStr string) Level {
	switch levelStr {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Debug 调试日志
func Debug(msg string, fields ...Field) {
	if globalLogger != nil {
		globalLogger.Debug(msg, fields...)
	}
}

// Info 信息日志
func Info(msg string, fields ...Field) {
	if globalLogger != nil {
		globalLogger.Info(msg, fields...)
	}
}

// Warn 警告日志
func Warn(msg string, fields ...Field) {
	if globalLogger != nil {
		globalLogger.Warn(msg, fields...)
	}
}

// Error 错误日志
func Error(msg string, fields ...Field) {
	if globalLogger != nil {
		globalLogger.Error(msg, fields...)
	}
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...Field) {
	if globalLogger != nil {
		globalLogger.Fatal(msg, fields...)
	}
	os.Exit(1)
}

// 实现 simpleLogger 接口
func (l *simpleLogger) Debug(msg string, fields ...Field) {
	if l.level <= DebugLevel {
		l.log(DebugLevel, msg, fields...)
	}
}

func (l *simpleLogger) Info(msg string, fields ...Field) {
	if l.level <= InfoLevel {
		l.log(InfoLevel, msg, fields...)
	}
}

func (l *simpleLogger) Warn(msg string, fields ...Field) {
	if l.level <= WarnLevel {
		l.log(WarnLevel, msg, fields...)
	}
}

func (l *simpleLogger) Error(msg string, fields ...Field) {
	if l.level <= ErrorLevel {
		l.log(ErrorLevel, msg, fields...)
	}
}

func (l *simpleLogger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields...)
}

func (l *simpleLogger) log(level Level, msg string, fields ...Field) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelName := levelNames[level]
	
	logMsg := timestamp + " [" + levelName + "] " + msg
	
	// 添加字段信息
	if len(fields) > 0 {
		logMsg += " {"
		for i, field := range fields {
			if i > 0 {
				logMsg += ", "
			}
			logMsg += field.Key + ": " + formatValue(field.Value)
		}
		logMsg += "}"
	}
	
	logMsg += "\n"
	l.writer.Write([]byte(logMsg))
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return "\"" + val + "\""
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%v", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", val)
	case float32, float64:
		return fmt.Sprintf("%v", val)
	case bool:
		return fmt.Sprintf("%v", val)
	case error:
		return "\"" + val.Error() + "\""
	default:
		return fmt.Sprintf("\"%v\"", val)
	}
}

// 便捷函数创建字段
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Uint(key string, value uint) Field {
	return Field{Key: key, Value: value}
}

func Err(key string, err error) Field {
	return Field{Key: key, Value: err}
}