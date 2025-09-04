-- 删除旧数据库（如果存在）
DROP DATABASE IF EXISTS `alert_system`;

-- 创建新数据库
CREATE DATABASE `alert_system` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用新数据库
USE `alert_system`;

-- 创建新的告警信息表（只包含message和recipient字段）
CREATE TABLE alerts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message TEXT NOT NULL COMMENT '告警信息',
    recipient VARCHAR(255) NOT NULL COMMENT '收件人',
    alert_time DATETIME NOT NULL COMMENT '告警时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_alert_time (alert_time),
    INDEX idx_recipient (recipient)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警信息表';

-- 插入一些测试数据
INSERT INTO alerts (message, recipient, alert_time) VALUES
('检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查', 'felixgao', '2025-09-04 15:11:11'),
('检测到域名【api.example.com】服务响应时间超过阈值，当前响应时间2.5秒', 'felixgao', '2025-09-04 15:26:11'),
('检测到域名【cdn.kugou.com】CDN节点异常，影响用户访问', 'felixgao', '2025-09-04 15:21:11'); 