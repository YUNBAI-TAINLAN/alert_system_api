package main

import (
	"time"
)

// Alert 告警信息结构
type Alert struct {
	ID          int       `json:"id" db:"id"`
	Domain      string    `json:"domain" db:"domain"`
	Message     string    `json:"message" db:"message"`
	Source      string    `json:"source" db:"source"`
	Status      string    `json:"status" db:"status"`
	Region      string    `json:"region" db:"region"`
	AlertTime   time.Time `json:"alert_time" db:"alert_time"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateAlertRequest 创建告警请求结构
type CreateAlertRequest struct {
	Domain    string `json:"domain" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Source    string `json:"source" binding:"required"`
	Status    string `json:"status"`
	Region    string `json:"region"`
	AlertTime string `json:"alert_time"`
}

// AlertResponse 告警响应结构
type AlertResponse struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    []Alert `json:"data,omitempty"`
}

// PeriodRequest 时间段查询请求
type PeriodRequest struct {
	StartTime string `json:"start_time" form:"start_time"`
	EndTime   string `json:"end_time" form:"end_time"`
} 