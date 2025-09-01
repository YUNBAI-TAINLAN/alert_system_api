package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error
	
	// 先连接到MySQL服务器（不指定数据库）
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port)
	
	// 连接MySQL服务器
	tempDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("连接MySQL服务器失败: %v", err)
	}
	defer tempDB.Close()
	
	// 测试连接
	if err = tempDB.Ping(); err != nil {
		return fmt.Errorf("MySQL服务器连接测试失败: %v", err)
	}
	
	// 创建数据库（如果不存在）
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", config.Database.Database)
	if _, err = tempDB.Exec(createDBSQL); err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}
	
	log.Printf("数据库 %s 创建/确认成功", config.Database.Database)
	
	// 现在连接到指定数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Database)
	
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}
	
	// 测试连接
	if err = db.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}
	
	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// 创建表
	if err = createTable(); err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}
	
	log.Println("数据库连接成功")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// createTable 创建告警信息表
func createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS alerts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		domain VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		source VARCHAR(100) NOT NULL,
		status VARCHAR(50),
		region VARCHAR(50),
		alert_time DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_alert_time (alert_time),
		INDEX idx_domain (domain),
		INDEX idx_source (source)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	
	_, err := db.Exec(query)
	return err
}

// InsertAlert 插入告警信息
func InsertAlert(alert *Alert) error {
	query := `
	INSERT INTO alerts (domain, message, source, status, region, alert_time)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	
	result, err := db.Exec(query, alert.Domain, alert.Message, alert.Source, alert.Status, alert.Region, alert.AlertTime)
	if err != nil {
		return fmt.Errorf("插入告警信息失败: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取插入ID失败: %v", err)
	}
	
	alert.ID = int(id)
	return nil
}

// GetAlerts 获取所有告警信息
func GetAlerts() ([]Alert, error) {
	query := `SELECT id, domain, message, source, status, region, alert_time, created_at, updated_at FROM alerts ORDER BY alert_time DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询告警信息失败: %v", err)
	}
	defer rows.Close()
	
	var alerts []Alert
	for rows.Next() {
		var alert Alert
		err := rows.Scan(&alert.ID, &alert.Domain, &alert.Message, &alert.Source, &alert.Status, &alert.Region, &alert.AlertTime, &alert.CreatedAt, &alert.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描告警信息失败: %v", err)
		}
		alerts = append(alerts, alert)
	}
	
	return alerts, nil
}

// GetAlertsByTimeRange 根据时间范围获取告警信息
func GetAlertsByTimeRange(startTime, endTime time.Time) ([]Alert, error) {
	query := `
	SELECT id, domain, message, source, status, region, alert_time, created_at, updated_at 
	FROM alerts 
	WHERE alert_time BETWEEN ? AND ? 
	ORDER BY alert_time DESC
	`
	
	rows, err := db.Query(query, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("查询时间段告警信息失败: %v", err)
	}
	defer rows.Close()
	
	var alerts []Alert
	for rows.Next() {
		var alert Alert
		err := rows.Scan(&alert.ID, &alert.Domain, &alert.Message, &alert.Source, &alert.Status, &alert.Region, &alert.AlertTime, &alert.CreatedAt, &alert.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描告警信息失败: %v", err)
		}
		alerts = append(alerts, alert)
	}
	
	return alerts, nil
} 