# 预警系统 API

这是一个用Go语言开发的预警信息管理系统，支持接收预警信息、定时统计和邮件通知功能。

## 功能特性

- 🔔 **预警信息接收**：通过HTTP API接口接收和存储预警信息
- ⏰ **定时统计**：每天晚上10点自动统计当天晚上7点到10点的预警信息
- 📧 **邮件通知**：将统计结果通过邮件发送给相关人员
- 🗄️ **数据持久化**：使用MySQL数据库存储预警信息
- 📊 **查询统计**：支持按时间段查询预警信息

## 系统架构

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  外部系统   │───▶│  HTTP API   │───▶│   MySQL     │
│  (预警源)   │    │  (Gin)      │    │  数据库     │
└─────────────┘    └─────────────┘    └─────────────┘
                          │
                          ▼
                   ┌─────────────┐
                   │  定时任务    │
                   │  (Cron)     │
                   └─────────────┘
                          │
                          ▼
                   ┌─────────────┐
                   │  邮件服务    │
                   │  (HTTP API)  │
                   └─────────────┘
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL 5.7+
- 支持SMTP的邮箱服务

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境变量

复制配置文件示例：
```bash
cp config.example .env
```

编辑 `.env` 文件，配置以下参数：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_DATABASE=alert_system

# 邮件配置 - 使用HTTP API方式
EMAIL_API_URL=http://opi.kgidc.cn/mail/email/send_email.php
EMAIL_APP_ID=v1-5f4769fe10c9c
EMAIL_APP_SECRET=c1e271982a82e325ef8ab5b0313fd102
EMAIL_FROM=system@company.com
EMAIL_TO=felixgao@kugou.net
EMAIL_DEBUG_MODE=false
EMAIL_DEBUG_API_URL=http://10.16.2.146:6709/mail/email/send_email.php

# 服务器配置
# 开发环境: localhost (只允许本机访问)
# 生产环境: 0.0.0.0 (允许外部访问)
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

### 4. 创建数据库

```sql
CREATE DATABASE `alert-api` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 运行程序

```bash
go run .
```

程序启动后会在指定端口（默认8080）监听HTTP请求。

## API 接口文档

### 1. 健康检查

**GET** `/health`

检查服务是否正常运行。

**响应示例：**
```json
{
  "status": "ok",
  "message": "服务正常运行"
}
```

### 2. 创建告警信息

**POST** `/api/v1/alerts`

**请求体：**
```json
{
  "domain": "example.com",
  "message": "服务器CPU使用率过高",
  "source": "监控系统",
  "status": "critical",
  "region": "北京",
  "alert_time": "2024-01-15 19:30:00"
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "告警信息创建成功",
  "data": {
    "id": 1,
    "domain": "example.com",
    "message": "服务器CPU使用率过高",
    "source": "监控系统",
    "status": "critical",
    "region": "北京",
    "alert_time": "2024-01-15T19:30:00Z",
    "created_at": "2024-01-15T19:30:00Z",
    "updated_at": "2024-01-15T19:30:00Z"
  }
}
```

### 3. 获取告警信息

**GET** `/api/v1/alerts?page=1&page_size=20`

**查询参数：**
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20，最大100）

**响应示例：**
```json
{
  "code": 200,
  "message": "获取告警信息成功",
  "data": [...],
  "total": 50,
  "page": 1,
  "size": 20
}
```

### 4. 按时间段查询告警

**GET** `/api/v1/alerts/period?start_time=2024-01-15 19:00:00&end_time=2024-01-15 22:00:00`

**查询参数：**
- `start_time`: 开始时间（格式：YYYY-MM-DD HH:mm:ss）
- `end_time`: 结束时间（格式：YYYY-MM-DD HH:mm:ss）

**响应示例：**
```json
{
  "code": 200,
  "message": "获取告警信息成功",
  "data": [...],
  "start_time": "2024-01-15 19:00:00",
  "end_time": "2024-01-15 22:00:00",
  "total": 15
}
```

## 定时任务

系统会在每天晚上10点自动执行以下任务：

1. 统计当天晚上7点到10点的告警信息
2. 生成统计报告
3. 通过邮件发送给配置的收件人

## 数据库结构

### alerts 表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | INT | 主键，自增 |
| domain | VARCHAR(255) | 域名 |
| message | TEXT | 告警信息 |
| source | VARCHAR(100) | 告警来源 |
| status | VARCHAR(50) | 告警状态 |
| region | VARCHAR(50) | 区域 |
| alert_time | DATETIME | 告警时间 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

## 邮件模板

系统使用HTML模板生成邮件内容，包含：

- 统计摘要（总告警数量、涉及域名数量、告警来源数量）
- 详细告警信息表格
- 时间范围信息

### 常见问题

1. **数据库连接失败**
   - 检查数据库配置
   - 确认数据库服务运行状态
   - 验证用户名密码

2. **邮件发送失败**
   - 检查SMTP配置
   - 确认邮箱密码（可能需要应用专用密码）
   - 验证网络连接

3. **定时任务不执行**
   - 检查系统时间
   - 查看程序日志
   - 确认cron表达式

## 开发说明

### 项目结构

```
api/
├── main.go          # 程序入口
├── config.go        # 配置管理
├── models.go        # 数据模型
├── handlers.go      # HTTP处理器
├── database.go      # 数据库操作
├── email.go         # 邮件服务
├── go.mod           # Go模块文件
├── go.sum           # 依赖校验
├── config.example   # 配置示例
├── init.sql         # 数据库初始化脚本
└── README.md        # 项目文档
```

### 扩展功能建议

1. **告警级别管理**：支持不同级别的告警处理
2. **告警去重**：避免重复告警
3. **告警确认机制**：支持告警确认和处理状态跟踪
4. **Web界面**：提供Web管理界面
5. **告警规则配置**：支持自定义告警规则
6. **多租户支持**：支持多组织使用

## 许可证

MIT License 