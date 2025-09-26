package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

// Config 日志配置
type Config struct {
	Level      string `json:"level" yaml:"level"`             // 日志级别
	Format     string `json:"format" yaml:"format"`           // 日志格式 json/text
	Output     string `json:"output" yaml:"output"`           // 输出方式 file/console/both
	FilePath   string `json:"file_path" yaml:"file_path"`     // 日志文件路径
	MaxSize    int    `json:"max_size" yaml:"max_size"`       // 单个文件最大尺寸(MB)
	MaxBackups int    `json:"max_backups" yaml:"max_backups"` // 保留旧文件的最大个数
	MaxAge     int    `json:"max_age" yaml:"max_age"`         // 保留旧文件的最大天数
	Compress   bool   `json:"compress" yaml:"compress"`       // 是否压缩
}

// DefaultConfig 默认日志配置
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		Format:     "json",
		Output:     "both",
		FilePath:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}
}

// Init 初始化日志系统
func Init(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	Logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 设置日志格式
	if config.Format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 设置输出
	var writers []io.Writer

	switch config.Output {
	case "file":
		writers = append(writers, getFileWriter(config))
	case "console":
		writers = append(writers, os.Stdout)
	case "both":
		writers = append(writers, os.Stdout, getFileWriter(config))
	default:
		writers = append(writers, os.Stdout)
	}

	Logger.SetOutput(io.MultiWriter(writers...))

	Logger.Info("日志系统初始化完成")
	return nil
}

// getFileWriter 获取文件写入器
func getFileWriter(config *Config) io.Writer {
	// 确保日志目录存在
	logDir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
		return os.Stdout
	}

	return &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
}

// GinLogger 返回gin的日志中间件
func GinLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf(`{"time":"%s","method":"%s","path":"%s","protocol":"%s","status_code":%d,"latency":"%s","client_ip":"%s","user_agent":"%s","error_message":"%s"}%s`,
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Request.UserAgent(),
				param.ErrorMessage,
				"\n",
			)
		},
		Output: Logger.Out,
	})
}

// GinRecovery 返回gin的恢复中间件
func GinRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		Logger.WithFields(logrus.Fields{
			"error":  recovered,
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"ip":     c.ClientIP(),
		}).Error("服务器内部错误")

		c.JSON(500, gin.H{
			"error":   "internal_server_error",
			"message": "服务器内部错误",
			"code":    500,
		})
	})
}

// Info 记录信息日志
func Info(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Info(msg)
	} else {
		Logger.Info(msg)
	}
}

// Error 记录错误日志
func Error(msg string, err error, fields ...logrus.Fields) {
	entry := Logger.WithField("error", err)
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Error(msg)
}

// Warn 记录警告日志
func Warn(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Warn(msg)
	} else {
		Logger.Warn(msg)
	}
}

// Debug 记录调试日志
func Debug(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Debug(msg)
	} else {
		Logger.Debug(msg)
	}
}
