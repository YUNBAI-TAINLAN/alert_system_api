package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level"`      // 日志级别: debug, info, warn, error
	FilePath   string `json:"file_path"`  // 日志文件路径
	MaxSize    int    `json:"max_size"`   // 单个日志文件最大大小(MB)
	MaxBackups int    `json:"max_backups"` // 保留的日志文件数量
	MaxAge     int    `json:"max_age"`    // 日志文件保留天数
	Compress   bool   `json:"compress"`   // 是否压缩旧日志文件
	Console    bool   `json:"console"`    // 是否同时输出到控制台
}

// InitLogger 初始化日志系统
func InitLogger(config LogConfig) error {
	Logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 设置日志格式为JSON格式，便于解析
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 创建日志目录
	logDir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 配置日志轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,    // MB
		MaxBackups: config.MaxBackups, // 保留文件数
		MaxAge:     config.MaxAge,     // 天数
		Compress:   config.Compress,   // 压缩
	}

	// 设置输出目标
	if config.Console {
		// 同时输出到文件和控制台
		multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
		Logger.SetOutput(multiWriter)
	} else {
		// 只输出到文件
		Logger.SetOutput(lumberjackLogger)
	}

	Logger.WithFields(logrus.Fields{
		"level":       config.Level,
		"file_path":   config.FilePath,
		"max_size":    config.MaxSize,
		"max_backups": config.MaxBackups,
		"max_age":     config.MaxAge,
		"compress":    config.Compress,
		"console":     config.Console,
	}).Info("日志系统初始化成功")

	return nil
}

// LogRequest 记录HTTP请求日志
func LogRequest(method, path, clientIP, userAgent string, statusCode int, latency string) {
	Logger.WithFields(logrus.Fields{
		"type":        "request",
		"method":      method,
		"path":        path,
		"client_ip":   clientIP,
		"user_agent":  userAgent,
		"status_code": statusCode,
		"latency":     latency,
	}).Info("HTTP请求处理完成")
}

// LogEmail 记录邮件发送日志
func LogEmail(recipient, subject string, success bool, errorMsg string) {
	fields := logrus.Fields{
		"type":      "email",
		"recipient": recipient,
		"subject":   subject,
		"success":   success,
	}
	
	if errorMsg != "" {
		fields["error"] = errorMsg
		Logger.WithFields(fields).Error("邮件发送失败")
	} else {
		Logger.WithFields(fields).Info("邮件发送成功")
	}
}

// LogDatabase 记录数据库操作日志
func LogDatabase(operation, table string, success bool, errorMsg string, rowsAffected int64) {
	fields := logrus.Fields{
		"type":      "database",
		"operation": operation,
		"table":     table,
		"success":   success,
	}
	
	if success {
		fields["rows_affected"] = rowsAffected
		Logger.WithFields(fields).Info("数据库操作成功")
	} else {
		if errorMsg != "" {
			fields["error"] = errorMsg
		}
		Logger.WithFields(fields).Error("数据库操作失败")
	}
}

// LogCronJob 记录定时任务日志
func LogCronJob(jobName string, success bool, message string, duration string) {
	fields := logrus.Fields{
		"type":     "cron",
		"job_name": jobName,
		"success":  success,
		"duration": duration,
	}
	
	if success {
		Logger.WithFields(fields).Info(message)
	} else {
		Logger.WithFields(fields).Error(message)
	}
}

// LogAlert 记录告警相关日志
func LogAlert(operation string, alertID int64, recipient, message string, success bool, errorMsg string) {
	fields := logrus.Fields{
		"type":      "alert",
		"operation": operation,
		"alert_id":  alertID,
		"recipient": recipient,
		"message":   message,
		"success":   success,
	}
	
	if errorMsg != "" {
		fields["error"] = errorMsg
		Logger.WithFields(fields).Error("告警操作失败")
	} else {
		Logger.WithFields(fields).Info("告警操作成功")
	}
}

// LogSystem 记录系统级日志
func LogSystem(level logrus.Level, component, message string, fields map[string]interface{}) {
	logFields := logrus.Fields{
		"type":      "system",
		"component": component,
	}
	
	// 合并额外字段
	for k, v := range fields {
		logFields[k] = v
	}
	
	Logger.WithFields(logFields).Log(level, message)
} 