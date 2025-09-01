-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `alert-api` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE `alert-api`;

-- 创建告警信息表
CREATE TABLE IF NOT EXISTS alerts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    domain VARCHAR(255) NOT NULL COMMENT '域名',
    message TEXT NOT NULL COMMENT '告警信息',
    source VARCHAR(100) NOT NULL COMMENT '告警来源',
    status VARCHAR(50) COMMENT '告警状态',
    region VARCHAR(50) COMMENT '区域',
    alert_time DATETIME NOT NULL COMMENT '告警时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_alert_time (alert_time),
    INDEX idx_domain (domain),
    INDEX idx_source (source),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警信息表';

-- 插入一些示例数据
INSERT INTO alerts (domain, message, source, status, region, alert_time) VALUES
('search.suggest.kgidc.cn', '检测到域名【search.suggest.kgidc.cn】北方已切量,但南方超过24小时未切量,请检查', 'RPC后台', 'warning', 'north', '2025-08-28 19:30:00'),
('api.example.com', '域名【api.example.com】响应时间超过5秒', '监控系统', 'error', 'south', '2025-08-28 20:15:00'),
('web.example.com', '域名【web.example.com】连接数达到上限', '负载均衡器', 'critical', 'east', '2025-08-28 20:45:00'); 