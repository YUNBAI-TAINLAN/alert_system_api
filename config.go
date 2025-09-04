package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Config 配置结构
type Config struct {
	Database DatabaseConfig
	Email    EmailConfig
	Server   ServerConfig
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string
	Port string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载.env文件
	loadEnvFile()
	
	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			Database: getEnv("DB_DATABASE", "alert_message"),
		},
		Email: EmailConfig{
			APIUrl:      getEnv("EMAIL_API_URL", "http://opi.kgidc.cn/mail/email/send_email.php"),
			AppID:       getEnv("EMAIL_APP_ID", "v1-5f4769fe10c9c"),
			AppSecret:   getEnv("EMAIL_APP_SECRET", "c1e271982a82e325ef8ab5b0313fd102"),
			From:        getEnv("EMAIL_FROM", "system@company.com"),
			To:          []string{}, // 不再使用固定收件人列表
			DebugMode:   getEnvAsBool("EMAIL_DEBUG_MODE", false),
			DebugAPIUrl: getEnv("EMAIL_DEBUG_API_URL", "http://10.16.2.146:6709/mail/email/send_email.php"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
	
	return config
}

// loadEnvFile 加载.env文件
func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		// .env文件不存在，使用环境变量
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为布尔值
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsSlice 获取环境变量并转换为字符串切片
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// 支持逗号分隔的多个值
		parts := strings.Split(value, ",")
		var result []string
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}