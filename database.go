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
		message TEXT NOT NULL,
		recipient VARCHAR(255) NOT NULL,
		alert_time DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_alert_time (alert_time),
		INDEX idx_recipient (recipient)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	
	_, err := db.Exec(query)
	return err
}

// InsertAlert 插入告警信息
func InsertAlert(alert *Alert) error {
	query := `
	INSERT INTO alerts (message, recipient, alert_time)
	VALUES (?, ?, ?)
	`
	
	result, err := db.Exec(query, alert.Message, alert.Recipient, alert.AlertTime)
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
	query := `SELECT id, message, recipient, alert_time, created_at, updated_at FROM alerts ORDER BY alert_time DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询告警信息失败: %v", err)
	}
	defer rows.Close()
	
	var alerts []Alert
	for rows.Next() {
		var alert Alert
		err := rows.Scan(&alert.ID, &alert.Message, &alert.Recipient, &alert.AlertTime, &alert.CreatedAt, &alert.UpdatedAt)
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
	SELECT id, message, recipient, alert_time, created_at, updated_at 
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
		err := rows.Scan(&alert.ID, &alert.Message, &alert.Recipient, &alert.AlertTime, &alert.CreatedAt, &alert.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描告警信息失败: %v", err)
		}
		alerts = append(alerts, alert)
	}
	
	return alerts, nil
}

// GetAlertsByRecipient 根据收件人获取告警信息
func GetAlertsByRecipient(recipient string) ([]Alert, error) {
	query := `
	SELECT id, message, recipient, alert_time, created_at, updated_at 
	FROM alerts 
	WHERE recipient = ? 
	ORDER BY alert_time DESC
	`
	
	rows, err := db.Query(query, recipient)
	if err != nil {
		return nil, fmt.Errorf("查询收件人告警信息失败: %v", err)
	}
	defer rows.Close()
	
	var alerts []Alert
	for rows.Next() {
		var alert Alert
		err := rows.Scan(&alert.ID, &alert.Message, &alert.Recipient, &alert.AlertTime, &alert.CreatedAt, &alert.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描告警信息失败: %v", err)
		}
		alerts = append(alerts, alert)
	}
	
	return alerts, nil
}

// GetAlertsByTimeRangeAndRecipient 根据时间范围和收件人获取告警信息
func GetAlertsByTimeRangeAndRecipient(startTime, endTime time.Time, recipient string) ([]Alert, error) {
	query := `
	SELECT id, message, recipient, alert_time, created_at, updated_at 
	FROM alerts 
	WHERE alert_time BETWEEN ? AND ? AND recipient = ?
	ORDER BY alert_time DESC
	`
	
	rows, err := db.Query(query, startTime, endTime, recipient)
	if err != nil {
		return nil, fmt.Errorf("查询时间段和收件人告警信息失败: %v", err)
	}
	defer rows.Close()
	
	var alerts []Alert
	for rows.Next() {
		var alert Alert
		err := rows.Scan(&alert.ID, &alert.Message, &alert.Recipient, &alert.AlertTime, &alert.CreatedAt, &alert.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描告警信息失败: %v", err)
		}
		alerts = append(alerts, alert)
	}
	
	return alerts, nil
} 

// GetUniqueRecipients 获取所有唯一的收件人
func GetUniqueRecipients() ([]string, error) {
	query := `SELECT DISTINCT recipient FROM alerts ORDER BY recipient`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询唯一收件人失败: %v", err)
	}
	defer rows.Close()
	
	var recipients []string
	for rows.Next() {
		var recipient string
		err := rows.Scan(&recipient)
		if err != nil {
			return nil, fmt.Errorf("扫描收件人失败: %v", err)
		}
		recipients = append(recipients, recipient)
	}
	
	return recipients, nil
}

// GetAlertsGroupedByRecipient 根据时间范围获取按收件人分组的告警信息
func GetAlertsGroupedByRecipient(startTime, endTime time.Time) ([]UserAlerts, error) {
	// 先获取所有唯一的收件人
	recipients, err := GetUniqueRecipients()
	if err != nil {
		return nil, err
	}
	
	var userAlertsList []UserAlerts
	
	for _, recipient := range recipients {
		alerts, err := GetAlertsByTimeRangeAndRecipient(startTime, endTime, recipient)
		if err != nil {
			return nil, err
		}
		
		if len(alerts) > 0 {
			userAlertsList = append(userAlertsList, UserAlerts{
				Recipient: recipient,
				Alerts:    alerts,
			})
		}
	}
	
	return userAlertsList, nil
}