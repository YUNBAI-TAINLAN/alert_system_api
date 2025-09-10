# 预警系统 API

这是一个用Go语言开发的预警信息管理系统，支持接收预警信息、定时统计和邮件通知功能。

## 功能特性

- 🔔 **预警信息接收**：通过HTTP API接口接收和存储预警信息
- ⏰ **定时统计**：每天晚上10点自动统计当天晚上7点到10点的预警信息
- 📧 **邮件通知**：将统计结果通过邮件发送给相关人员
- 🗄️ **数据持久化**：使用MySQL数据库存储预警信息
- 📊 **查询统计**：支持按时间段查询预警信息
- 👥 **用户分组**：支持按收件人分组发送告警信息
- 🔗 **动态收件人**：根据告警信息中的收件人变量自动生成邮箱地址
- 📋 **用户列表管理**：从userlist.json文件加载用户信息，支持英文名到邮箱地址的映射
- 🎨 **管理员邮件**：当用户未找到时，自动发送合并邮件给管理员
- 🎨 **美观邮件模板**：HTML格式邮件，支持中文显示，不显示技术性ID信息

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

## 数据模型

### 告警信息结构（简化版）

```json
{
  "id": 1,
  "message": "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
  "recipient": "zhangsan",
  "alert_time": "2025-01-15 19:30:00",
  "created_at": "2025-01-15 19:30:00",
  "updated_at": "2025-01-15 19:30:00"
}
```

**字段说明：**
- `message`: 告警信息内容（必填）
- `recipient`: 收件人标识（必填，系统会自动添加@kugou.net后缀）
- `alert_time`: 告警时间（可选，默认为当前时间）

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

编辑 `.env` 文件：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_DATABASE=alert_system

# 邮件配置
EMAIL_API_URL=http://opi.kgidc.cn/mail/email/send_email.php
EMAIL_APP_ID=v1-5f4769fe10c9c
EMAIL_APP_SECRET=c1e271982a82e325ef8ab5b0313fd102
EMAIL_FROM=system@company.com
EMAIL_TO=admin@company.com
EMAIL_DEBUG_MODE=false

# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

### 4. 数据库初始化

#### 新安装
```bash
mysql -u root -p < init.sql
```

#### 从旧版本升级
```bash
mysql -u root -p < migration.sql
```

### 5. 启动服务

```bash
go run .
```

服务将在 `http://localhost:8080` 启动。

## API 接口

### 1. 创建告警信息

```bash
POST /api/v1/alerts
Content-Type: application/json

{
  "message": "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
  "recipient": "zhangsan",
  "alert_time": "2025-01-15 19:30:00"
}
```

### 2. 获取告警信息

```bash
GET /api/v1/alerts?page=1&page_size=20
```

### 3. 按收件人查询告警

```bash
GET /api/v1/alerts/recipient?recipient=zhangsan
```

### 4. 按时间段查询告警

```bash
GET /api/v1/alerts/period?start_time=2025-01-15 19:00:00&end_time=2025-01-15 22:00:00
```

### 5. 测试邮件发送

```bash
POST /test-email
```

## 定时任务

系统每天晚上10点自动执行以下任务：

1. 获取当天晚上7点到10点的告警信息
2. 按收件人分组
3. 为每个收件人发送单独的邮件通知

## 邮件功能

### 动态收件人

系统会根据告警信息中的 `recipient` 字段自动生成邮箱地址：

1. **用户列表匹配**：首先在 `userlist.json` 中查找对应的英文名
2. **邮箱地址生成**：
   - 如果 `recipient` 已经包含 `@` 符号，直接使用
   - 否则自动添加 `@kugou.net` 后缀
3. **管理员邮件**：当用户未找到时，发送合并邮件给管理员（liyongchang@kugou.net）

例如：
- `recipient: "zhangsan"` → 邮箱：`zhangsan@kugou.net`
- `recipient: "lisi@company.com"` → 邮箱：`lisi@company.com`
- `recipient: "unknownuser"` → 发送给管理员：`liyongchang@kugou.net`

### 邮件模板

每个用户会收到包含其专属告警信息的邮件，邮件内容包括：

- 收件人信息
- 告警统计摘要
- 详细的告警信息列表
- 美观的HTML格式

#### 管理员邮件
当用户未找到时，管理员（liyongchang@kugou.net）会收到包含所有未找到用户告警信息的合并邮件，内容包括：

- 未找到用户列表
- 所有相关告警信息
- 处理建议

## 测试

### 运行测试脚本

```bash
# 给脚本执行权限
chmod +x test_new_api.sh

# 运行测试
./test_new_api.sh
```

### 手动测试

1. 创建告警信息：
```bash
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "message": "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
    "recipient": "zhangsan"
  }'
```

2. 测试邮件发送：
```bash
curl -X POST http://localhost:8080/test-email
```

## 版本更新

### v2.0.0 更新内容

1. **简化数据模型**：只保留 `message` 和 `recipient` 字段
2. **动态收件人**：根据告警信息自动生成邮箱地址
3. **用户分组**：按收件人分组发送告警信息
4. **优化邮件模板**：为每个用户生成专属邮件

### 从旧版本升级

1. 备份现有数据
2. 运行数据库迁移脚本：`mysql -u root -p < migration.sql`
3. 更新代码并重启服务

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库配置
   - 确保MySQL服务正在运行
   - 验证用户名和密码

2. **邮件发送失败**
   - 检查邮件API配置
   - 验证App ID和App Secret
   - 查看日志中的详细错误信息

3. **定时任务不执行**
   - 检查系统时间
   - 查看日志中的定时任务状态
   - 验证Cron表达式

### 日志查看

服务日志会显示详细的执行信息，包括：

- 数据库操作
- 邮件发送状态
- 定时任务执行情况
- API请求处理

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 许可证

MIT License 