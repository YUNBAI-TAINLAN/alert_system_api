package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateAlert 创建预警信息
func CreateAlert(c *gin.Context) {
	var req CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 解析预警时间
	var alertTime time.Time
	var err error
	if req.AlertTime != "" {
		alertTime, err = time.Parse("2006-01-02 15:04:05", req.AlertTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "预警时间格式错误，请使用 YYYY-MM-DD HH:mm:ss 格式",
			})
			return
		}
	} else {
		alertTime = time.Now()
	}

	// 创建预警对象
	alert := &Alert{
		Domain:    req.Domain,
		Message:   req.Message,
		Source:    req.Source,
		Status:    req.Status,
		Region:    req.Region,
		AlertTime: alertTime,
	}

	// 插入数据库
	if err := InsertAlert(alert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "存储预警信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "预警信息创建成功",
		"data":    alert,
	})
}

// GetAlertsHandler 获取所有预警信息
func GetAlertsHandler(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	alerts, err := GetAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取预警信息失败: " + err.Error(),
		})
		return
	}

	// 简单的分页处理
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(alerts) {
		alerts = []Alert{}
	} else if end > len(alerts) {
		alerts = alerts[start:]
	} else {
		alerts = alerts[start:end]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取预警信息成功",
		"data":    alerts,
		"total":   len(alerts),
		"page":    page,
		"size":    pageSize,
	})
}

// GetAlertsByPeriod 根据时间段获取预警信息
func GetAlertsByPeriod(c *gin.Context) {
	var req PeriodRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 解析时间参数
	var startTime, endTime time.Time
	var err error

	if req.StartTime != "" {
		startTime, err = time.Parse("2006-01-02 15:04:05", req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "开始时间格式错误，请使用 YYYY-MM-DD HH:mm:ss 格式",
			})
			return
		}
	} else {
		// 默认查询今天的数据
		now := time.Now()
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 19, 0, 0, 0, now.Location())
	}

	if req.EndTime != "" {
		endTime, err = time.Parse("2006-01-02 15:04:05", req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "结束时间格式错误，请使用 YYYY-MM-DD HH:mm:ss 格式",
			})
			return
		}
	} else {
		// 默认查询到今天结束
		now := time.Now()
		endTime = time.Date(now.Year(), now.Month(), now.Day(), 22, 59, 59, 999999999, now.Location())
	}

	alerts, err := GetAlertsByTimeRange(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取预警信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"message":   "获取预警信息成功",
		"data":      alerts,
		"start_time": startTime.Format("2006-01-02 15:04:05"),
		"end_time":   endTime.Format("2006-01-02 15:04:05"),
		"total":     len(alerts),
	})
} 