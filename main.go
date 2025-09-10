package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var config *Config

func main() {
	// 加载配置
	config = LoadConfig()

	// 初始化日志系统
	if err := InitLogger(config.Log); err != nil {
		log.Fatal("日志系统初始化失败:", err)
	}
	
	LogSystem(logrus.InfoLevel, "main", "告警系统启动", map[string]interface{}{
		"version": "1.0.0",
	})

	// 初始化数据库连接
	if err := InitDB(); err != nil {
		LogSystem(logrus.FatalLevel, "main", "数据库初始化失败", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatal("数据库初始化失败:", err)
	}
	defer CloseDB()
	
	LogSystem(logrus.InfoLevel, "main", "数据库连接成功", nil)

	// 初始化邮件配置
	InitEmailConfig()
	LogSystem(logrus.InfoLevel, "main", "邮件配置初始化完成", nil)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin路由
	r := gin.Default()
	
	// 添加日志中间件
	r.Use(LoggerMiddleware())

	// 设置路由
	setupRoutes(r)

	// 启动定时任务
	startCronJob()

	// 启动HTTP服务器
	serverAddr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	LogSystem(logrus.InfoLevel, "main", "服务器启动", map[string]interface{}{
		"address": serverAddr,
	})
	log.Printf("服务器启动在 %s...", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		LogSystem(logrus.FatalLevel, "main", "服务器启动失败", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatal("服务器启动失败:", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "服务正常运行"})
	})

	// 配置检查接口
	r.GET("/config", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"email_config": gin.H{
				"api_url":        emailConfig.APIUrl,
				"app_id":         emailConfig.AppID,
				"app_secret":     "***hidden***",
				"from":           emailConfig.From,
				"debug_mode":     emailConfig.DebugMode,
				"debug_api_url":  emailConfig.DebugAPIUrl,
				"note":           "收件人现在根据告警信息动态生成",
			},
			"cron_config": gin.H{
				"enabled":      config.Cron.Enabled,
				"schedule":     config.Cron.Schedule,
				"start_time":   fmt.Sprintf("%02d:%02d", config.Cron.StartHour, config.Cron.StartMinute),
				"end_time":     fmt.Sprintf("%02d:%02d", config.Cron.EndHour, config.Cron.EndMinute),
				"description":  "定时任务配置信息",
			},
		})
	})

	// 邮件测试接口
	r.POST("/test-email", func(c *gin.Context) {
		// 显示当前邮件配置信息
		log.Printf("邮件配置信息:")
		log.Printf("  API地址: %s", emailConfig.APIUrl)
		log.Printf("  App ID: %s", emailConfig.AppID)
		log.Printf("  调试模式: %v", emailConfig.DebugMode)
		if emailConfig.DebugMode {
			log.Printf("  调试API地址: %s", emailConfig.DebugAPIUrl)
		}
		log.Printf("  收件人: 根据告警信息动态生成")

		// 创建测试预警数据（按用户分组），包含多个用户
		testUserAlerts := []UserAlerts{
			{
				Recipient: "felixgao",
				Alerts: []Alert{
					{
						ID:        1,
						Message:   "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
						Recipient: "felixgao",
						AlertTime: time.Now().Add(-30 * time.Minute),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        2,
						Message:   "检测到域名【api.example.com】服务响应时间超过阈值，当前响应时间2.5秒",
						Recipient: "felixgao",
						AlertTime: time.Now().Add(-15 * time.Minute),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				},
			},
			{
				Recipient: "hugoli",
				Alerts: []Alert{
					{
						ID:        3,
						Message:   "检测到域名【cdn.kugou.com】CDN节点异常，影响用户访问",
						Recipient: "hugoli",
						AlertTime: time.Now().Add(-20 * time.Minute),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        4,
						Message:   "检测到数据库连接池使用率超过80%，请检查数据库性能",
						Recipient: "hugoli",
						AlertTime: time.Now().Add(-10 * time.Minute),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				},
			},
			{
				Recipient: "zhangsan",
				Alerts: []Alert{
					{
						ID:        5,
						Message:   "检测到服务器CPU使用率超过90%，请检查系统负载",
						Recipient: "zhangsan",
						AlertTime: time.Now().Add(-25 * time.Minute),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				},
			},
			{
				Recipient: "lisi",
				Alerts: []Alert{
					{
						ID:        6,
						Message:   "检测到内存使用率超过85%，请检查内存泄漏",
						Recipient: "lisi",
						AlertTime: time.Now().Add(-5 * time.Minute),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        7,
						Message:   "检测到磁盘空间不足，剩余空间小于10%",
						Recipient: "lisi",
						AlertTime: time.Now().Add(-2 * time.Minute),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				},
			},
		}

		// 发送测试邮件
		if err := SendAlertEmail(testUserAlerts); err != nil {
			log.Printf("邮件发送失败: %v", err)
			c.JSON(500, gin.H{
				"code":    500,
				"message": "邮件发送失败: " + err.Error(),
			})
			return
		}

		// 计算总告警数量
		totalAlerts := 0
		for _, userAlerts := range testUserAlerts {
			totalAlerts += len(userAlerts.Alerts)
		}

		log.Printf("测试邮件发送成功，共发送给 %d 个用户，总计 %d 条告警信息", len(testUserAlerts), totalAlerts)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "测试邮件发送成功",
			"data": gin.H{
				"user_count":   len(testUserAlerts),
				"total_alerts": totalAlerts,
				"recipients":   []string{"felixgao", "hugoli", "zhangsan", "lisi"},
				"details": []gin.H{
					{"user": "felixgao", "alert_count": 2},
					{"user": "hugoli", "alert_count": 2},
					{"user": "zhangsan", "alert_count": 1},
					{"user": "lisi", "alert_count": 2},
				},
			},
		})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 存储预警信息
		api.POST("/alerts", CreateAlert)
		
		// 获取预警信息
		api.GET("/alerts", GetAlertsHandler)
		
		// 获取指定时间段的预警信息
		api.GET("/alerts/period", GetAlertsByPeriod)
		
		// 根据收件人获取预警信息
		api.GET("/alerts/recipient", GetAlertsByRecipientHandler)
	}
}

func startCronJob() {
	// 检查是否启用定时任务
	if !config.Cron.Enabled {
		LogSystem(logrus.InfoLevel, "cron", "定时任务已禁用", nil)
		log.Println("定时任务已禁用")
		return
	}

	c := cron.New(cron.WithLocation(time.Local))
	
	// 使用配置的cron表达式执行定时任务
	_, err := c.AddFunc(config.Cron.Schedule, func() {
		startTime := time.Now()
		LogCronJob("alert_notification", true, "定时任务开始执行", "")
		
		// 使用配置的时间范围获取预警信息
		now := time.Now()
		alertStartTime := time.Date(now.Year(), now.Month(), now.Day(), 
			config.Cron.StartHour, config.Cron.StartMinute, 0, 0, now.Location())
		alertEndTime := time.Date(now.Year(), now.Month(), now.Day(), 
			config.Cron.EndHour, config.Cron.EndMinute, 0, 0, now.Location())
		
		LogSystem(logrus.InfoLevel, "cron", "查询告警时间范围", map[string]interface{}{
			"start_time": alertStartTime.Format("2006-01-02 15:04:05"),
			"end_time": alertEndTime.Format("2006-01-02 15:04:05"),
			"schedule": config.Cron.Schedule,
		})
		
		// 按收件人分组获取告警信息
		userAlertsList, err := GetAlertsGroupedByRecipient(alertStartTime, alertEndTime)
		if err != nil {
			duration := time.Since(startTime).String()
			LogCronJob("alert_notification", false, "获取预警信息失败: "+err.Error(), duration)
			log.Printf("获取预警信息失败: %v", err)
			return
		}
		
		if len(userAlertsList) == 0 {
			duration := time.Since(startTime).String()
			LogCronJob("alert_notification", true, "指定时间段内没有预警信息", duration)
			log.Println("指定时间段内没有预警信息")
			return
		}
		
		LogSystem(logrus.InfoLevel, "cron", "准备发送邮件", map[string]interface{}{
			"user_count": len(userAlertsList),
		})
		
		// 按用户分组发送邮件
		if err := SendAlertEmail(userAlertsList); err != nil {
			duration := time.Since(startTime).String()
			LogCronJob("alert_notification", false, "发送邮件失败: "+err.Error(), duration)
			log.Printf("发送邮件失败: %v", err)
		} else {
			duration := time.Since(startTime).String()
			LogCronJob("alert_notification", true, fmt.Sprintf("成功发送预警通知邮件，涉及 %d 个用户", len(userAlertsList)), duration)
			log.Printf("成功发送预警通知邮件，涉及 %d 个用户", len(userAlertsList))
		}
	})
	
	if err != nil {
		LogSystem(logrus.FatalLevel, "cron", "添加定时任务失败", map[string]interface{}{
			"error": err.Error(),
			"schedule": config.Cron.Schedule,
		})
		log.Fatal("添加定时任务失败:", err)
	}
	
	c.Start()
	LogSystem(logrus.InfoLevel, "cron", "定时任务已启动", map[string]interface{}{
		"schedule": config.Cron.Schedule,
		"start_time": fmt.Sprintf("%02d:%02d", config.Cron.StartHour, config.Cron.StartMinute),
		"end_time": fmt.Sprintf("%02d:%02d", config.Cron.EndHour, config.Cron.EndMinute),
	})
	log.Printf("定时任务已启动，执行时间: %s，查询范围: %02d:%02d - %02d:%02d", 
		config.Cron.Schedule, 
		config.Cron.StartHour, config.Cron.StartMinute,
		config.Cron.EndHour, config.Cron.EndMinute)
}

// LoggerMiddleware Gin日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 处理请求
		c.Next()
		
		// 计算处理时间
		latency := time.Since(start)
		
		// 记录请求日志
		LogRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
			c.Writer.Status(),
			latency.String(),
		)
	}
} 

