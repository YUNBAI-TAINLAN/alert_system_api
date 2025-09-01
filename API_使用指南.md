# 预警系统 API 使用指南

## 📋 概述

这是一个用Go语言开发的预警信息管理系统API，支持接收预警信息、查询统计和邮件通知功能。

**服务地址**: `http://0.0.0.0:8080`
**API版本**: `v1`
**基础路径**: `/api/v1`

## 🔧 接口列表

### 1. 健康检查
检查服务是否正常运行

**请求**:
```bash
GET /health
```

### 2. 配置查看
查看当前系统配置信息

**请求**:
```bash
GET /config
```

**响应**:
```json
{
  "status": "ok",
  "email_config": {
    "api_url": "http://opi.kgidc.cn/mail/email/send_email.php",
    "app_id": "v1-5f4769fe10c9c",
    "app_secret": "***hidden***",
    "from": "system@company.com",
    "to": ["felixgao@kugou.net"],
    "debug_mode": false,
    "debug_api_url": "http://10.16.2.146:6709/mail/email/send_email.php"
  }
}
```

### 3. 邮件测试
测试邮件发送功能

**请求**:
```bash
POST /test-email
```

**响应**:
```json
{
  "code": 200,
  "message": "测试邮件发送成功",
  "data": {
    "alert_count": 2,
    "recipients": ["felixgao@kugou.net"]
  }
}
```

### 4. 创建预警信息

**响应**:
```json
{
  "status": "ok",
  "message": "服务正常运行"
}
```

### 5. 创建预警信息
向系统提交新的预警信息

**请求**:
```bash
POST /api/v1/alerts
Content-Type: application/json
```

**请求体**:
```json
{
  "domain": "search.suggest.kgidc.cn",
  "message": "北方已切量,但南方超过24小时未切量,请检查",
  "source": "RPC后台",
  "status": "active",
  "region": "全国",
  "alert_time": "2025-01-15 10:30:00"
}
```

**字段说明**:
- `domain` (必填): 域名
- `message` (必填): 预警信息内容
- `source` (必填): 预警来源
- `status` (可选): 预警状态，如 "active", "resolved", "warning"
- `region` (可选): 区域信息
- `alert_time` (可选): 预警时间，格式为 "YYYY-MM-DD HH:mm:ss"，默认为当前时间

**响应**:
```json
{
  "code": 200,
  "message": "预警信息创建成功",
  "data": {
    "id": 1,
    "domain": "search.suggest.kgidc.cn",
    "message": "北方已切量,但南方超过24小时未切量,请检查",
    "source": "RPC后台",
    "status": "active",
    "region": "全国",
    "alert_time": "2025-01-15T10:30:00Z",
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

### 6. 获取所有预警信息
获取系统中的所有预警信息，支持分页

**请求**:
```bash
GET /api/v1/alerts?page=1&size=20
```

**查询参数**:
- `page` (可选): 页码，默认为1
- `size` (可选): 每页数量，默认为20，最大100

**响应**:
```json
{
  "code": 200,
  "message": "获取预警信息成功",
  "data": [
    {
      "id": 1,
      "domain": "search.suggest.kgidc.cn",
      "message": "北方已切量,但南方超过24小时未切量,请检查",
      "source": "RPC后台",
      "status": "active",
      "region": "全国",
      "alert_time": "2025-01-15T10:30:00Z",
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:30:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "size": 20
}
```

### 7. 按时间段查询预警
根据时间范围查询预警信息

**请求**:
```bash
GET /api/v1/alerts/period?start_time=2025-01-15 07:00:00&end_time=2025-01-15 22:00:00
```

**查询参数**:
- `start_time` (可选): 开始时间，格式为 "YYYY-MM-DD HH:mm:ss"
- `end_time` (可选): 结束时间，格式为 "YYYY-MM-DD HH:mm:ss"

**注意**: 如果不提供时间参数，默认查询当天晚上7点到10点的数据

**响应**:
```json
{
  "code": 200,
  "message": "获取预警信息成功",
  "data": [
    {
      "id": 1,
      "domain": "search.suggest.kgidc.cn",
      "message": "北方已切量,但南方超过24小时未切量,请检查",
      "source": "RPC后台",
      "status": "active",
      "region": "全国",
      "alert_time": "2025-01-15T10:30:00Z",
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:30:00Z"
    }
  ],
  "start_time": "2025-01-15 07:00:00",
  "end_time": "2025-01-15 22:00:00",
  "total": 1
}
```

## 🚀 调用示例

### cURL 示例

#### 1. 创建预警信息
```bash
curl -X POST http://localhost:8080/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "search.suggest.kgidc.cn",
    "message": "北方已切量,但南方超过24小时未切量,请检查",
    "source": "RPC后台",
    "status": "active",
    "region": "全国"
  }'
```

#### 2. 获取所有预警
```bash
curl http://localhost:8080/api/v1/alerts
```

#### 3. 按时间段查询
```bash
curl "http://localhost:8080/api/v1/alerts/period?start_time=2025-01-15%2007:00:00&end_time=2025-01-15%2022:00:00"
```



## 📧 邮件通知

系统会在每天晚上10点自动统计当天晚上7点到10点的预警信息，并发送邮件通知。

邮件内容采用新的预警通知样式，包含：
- 统计摘要
- 详细的预警信息卡片
- 每个预警包含图标、标题、详细信息和时间戳

### 邮件配置

系统使用HTTP API方式发送邮件，需要在配置文件中设置以下参数：

```bash
# 邮件配置 - 使用HTTP API方式
EMAIL_API_URL=http://opi.kgidc.cn/mail/email/send_email.php
EMAIL_APP_ID=v1-5f4769fe10c9c
EMAIL_APP_SECRET=c1e271982a82e325ef8ab5b0313fd102
EMAIL_FROM=system@company.com
EMAIL_TO=felixgao@kugou.net
EMAIL_DEBUG_MODE=false
EMAIL_DEBUG_API_URL=http://10.16.2.146:6709/mail/email/send_email.php
```

**配置说明**：
- `EMAIL_API_URL`: 邮件发送API地址（生产环境）
- `EMAIL_APP_ID`: 应用ID
- `EMAIL_APP_SECRET`: 应用密钥
- `EMAIL_FROM`: 发件人地址
- `EMAIL_TO`: 收件人地址，多个地址用逗号分隔
- `EMAIL_DEBUG_MODE`: 是否启用调试模式（true/false）
- `EMAIL_DEBUG_API_URL`: 调试模式下的API地址

### 测试邮件发送

可以使用以下命令测试邮件发送功能：

```bash
# 运行邮件测试工具
go run test_email.go
```

## ⚠️ 注意事项

1. **时间格式**: 所有时间参数必须使用 "YYYY-MM-DD HH:mm:ss" 格式
2. **必填字段**: domain、message、source 为必填字段
3. **分页限制**: 每页最大返回100条记录
4. **服务状态**: 调用前建议先检查 `/health` 接口确认服务状态
5. **错误处理**: 所有接口都会返回标准的错误响应格式

## 🔍 错误响应格式

```json
{
  "code": 400,
  "message": "错误描述信息"
}
```

常见错误码：
- `400`: 请求参数错误
- `500`: 服务器内部错误

