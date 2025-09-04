-- 数据库迁移脚本：从旧结构迁移到新结构
-- 执行前请备份数据库！

USE `alert-api`;

-- 1. 创建新的表结构
CREATE TABLE IF NOT EXISTS alerts_new (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message TEXT NOT NULL COMMENT '告警信息',
    recipient VARCHAR(255) NOT NULL COMMENT '收件人',
    alert_time DATETIME NOT NULL COMMENT '告警时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_alert_time (alert_time),
    INDEX idx_recipient (recipient)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警信息表';

-- 2. 迁移现有数据（如果有的话）
-- 注意：这里需要根据实际情况调整收件人字段的值
INSERT INTO alerts_new (message, recipient, alert_time, created_at, updated_at)
SELECT 
    CONCAT('检测到域名【', domain, '】', message) as message,
    'default_user' as recipient,  -- 这里需要根据实际情况设置收件人
    alert_time,
    created_at,
    updated_at
FROM alerts;

-- 3. 删除旧表
DROP TABLE IF EXISTS alerts;

-- 4. 重命名新表
RENAME TABLE alerts_new TO alerts;

-- 5. 插入一些示例数据（新格式）
INSERT INTO alerts (message, recipient, alert_time) VALUES
('检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查', 'zhangsan', '2025-01-15 19:30:00'),
('检测到域名【api.example.com】响应时间超过5秒', 'lisi', '2025-01-15 20:15:00'),
('检测到域名【web.example.com】连接数达到上限', 'wangwu', '2025-01-15 20:45:00'); 