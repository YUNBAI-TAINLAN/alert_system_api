# 预警系统 API

一个基于Go语言开发的企业级预警信息管理系统，提供完整的告警信息接收、存储、统计和邮件通知功能。

## ✨ 核心特性

- 🔔 **高性能API**：基于Gin框架，支持高并发请求处理
- ⏰ **智能统计**：自动按时间段统计告警信息
- 📧 **智能邮件**：动态收件人生成，支持用户分组发送
- 🗄️ **数据持久化**：MySQL数据库存储，支持复杂查询
- ⏰ **定时任务**：Cron定时器，自动执行统计和邮件发送
- 👥 **用户管理**：支持用户列表管理，英文名到邮箱映射
- 🎨 **美观界面**：HTML邮件模板，支持中文显示
- 🔗 **灵活配置**：环境变量配置，支持调试模式

## 🏗️ 系统架构

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
<code_block_to_apply_changes_from>
alert_system_api/
├── main.go              # 主程序入口
├── config.go            # 配置管理
├── database.go          # 数据库操作
├── handlers.go          # API处理器
├── models.go            # 数据模型
├── email.go             # 邮件服务
├── logger.go            # 日志系统
├── init.sql             # 数据库初始化脚本
├── migration.sql        # 数据库迁移脚本
├── userlist.json        # 用户列表
├── config.example       # 配置文件示例
├── test_new_api.sh      # 测试脚本
└── README.md            # 项目文档
```

## 🎛️ 配置说明

### 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `DB_HOST` | 数据库主机 | localhost |
| `DB_PORT` | 数据库端口 | 3306 |
| `DB_USERNAME` | 数据库用户名 | root |
| `DB_PASSWORD` | 数据库密码 | - |
| `DB_DATABASE` | 数据库名称 | alert_system |
| `EMAIL_API_URL` | 邮件API地址 | - |
| `EMAIL_APP_ID` | 邮件服务App ID | - |
| `EMAIL_APP_SECRET` | 邮件服务App Secret | - |
| `SERVER_HOST` | 服务器监听地址 | 0.0.0.0 |
| `SERVER_PORT` | 服务器端口 | 8080 |

### 用户列表配置

`userlist.json` 文件格式：

```json
[
  {
    "name": "张三",
    "e_name": "zhangsan",
    "email": "zhangsan@kugou.net"
  }
]
```

## 🐛 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库服务是否运行
   - 验证连接配置和权限
   - 查看日志中的详细错误信息

2. **邮件发送失败**
   - 检查邮件API配置
   - 验证App ID和App Secret
   - 确认网络连接正常

3. **定时任务不执行**
   - 检查系统时间设置
   - 查看日志中的定时任务状态
   - 验证Cron表达式配置

### 日志查看

系统使用结构化日志，包含以下信息：

- 数据库操作记录
- 邮件发送状态
- 定时任务执行情况
- API请求处理日志

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 📄 许可证

MIT License

---

**版本**: v1.0.0  
**作者**: 架构一组
**最后更新**: 2025年9月11日

##  数据模型

### 告警信息结构

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
- `recipient`: 收件人标识（必填，支持以下格式：完整邮箱地址、英文名、或系统会自动在用户列表中查找对应邮箱）
- `alert_time`: 告警时间（可选，默认为当前时间）

## 🚀 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+
- 支持HTTP API的邮件服务

### 1. 克隆项目

```bash
git clone <repository-url>
cd alert_system_api
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境

复制配置文件：

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

### 4. 初始化数据库

```bash
mysql -u root -p < init.sql
```

### 5. 启动服务

```bash
go run .
```

服务将在 `http://localhost:8080` 启动。

## 📚 API 接口

### 基础接口

| 接口 | 方法 | 描述 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/config` | GET | 查看配置信息 |

### 告警管理

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/v1/alerts` | POST | 创建告警信息 |
| `/api/v1/alerts` | GET | 获取告警列表（分页） |
| `/api/v1/alerts/recipient` | GET | 按收件人查询 |
| `/api/v1/alerts/period` | GET | 按时间段查询 |

### 测试接口

| 接口 | 方法 | 描述 |
|------|------|------|
| `/test-email` | POST | 测试邮件发送功能 |

### 请求示例

#### 创建告警信息

```bash
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "message": "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
    "recipient": "zhangsan"
  }'
```

#### 获取告警列表

```bash
curl "http://localhost:8080/api/v1/alerts?page=1&page_size=20"
```

#### 按收件人查询

```bash
curl "http://localhost:8080/api/v1/alerts/recipient?recipient=zhangsan"
```

#### 按时间段查询

```bash
curl "http://localhost:8080/api/v1/alerts/period?start_time=2025-01-15%2019:00:00&end_time=2025-01-15%2022:00:00"
```

## ⏰ 定时任务

系统配置了以下定时任务：

- **执行时间**：每天晚上10点
- **统计范围**：当天晚上7点到10点的告警信息
- **处理流程**：
  1. 获取指定时间段的告警信息
  2. 按收件人分组
  3. 为每个收件人发送专属邮件
  4. 未找到用户发送给管理员

## 📧 邮件功能

### 动态收件人生成

系统智能处理收件人信息：

1. **完整邮箱地址**：如果recipient包含@符号，直接使用该邮箱地址
2. **用户列表匹配**：如果recipient不包含@符号，在 `userlist.json` 中查找对应英文名获取邮箱
3. **管理员兜底**：未找到用户时发送给 `liyongchang@kugou.net`

**处理示例：**
- `recipient: "zhangsan@company.com"` → 直接使用：`zhangsan@company.com`
- `recipient: "zhangsan"` → 查找用户列表：`zhangsan@kugou.net`（如果用户存在）
- `recipient: "unknownuser"` → 管理员兜底：`liyongchang@kugou.net`

### 邮件模板

- **用户邮件**：包含专属告警信息，美观的HTML格式
- **管理员邮件**：包含所有未找到用户的告警信息
- **中文支持**：完美支持中文显示，无乱码问题

## 🧪 测试

### 自动化测试

```bash
# 给脚本执行权限
chmod +x test_new_api.sh

# 运行测试
./test_new_api.sh
```

### 手动测试

```bash
# 测试邮件发送
curl -X POST http://localhost:8080/test-email

# 健康检查
curl http://localhost:8080/health
```

## 📁 项目结构

```
alert_system_api/
├── main.go              # 主程序入口
├── config.go            # 配置管理
├── database.go          # 数据库操作
├── handlers.go          # API处理器
├── models.go            # 数据模型
├── email.go             # 邮件服务
├── logger.go            # 日志系统
├── init.sql             # 数据库初始化脚本
├── migration.sql        # 数据库迁移脚本
├── userlist.json        # 用户列表
├── config.example       # 配置文件示例
├── test_new_api.sh      # 测试脚本
└── README.md            # 项目文档
``` 