# 预警系统 API 接口文档

## 1. 健康检查接口

**简要描述**: 检查预警系统服务是否正常运行，用于监控和健康检查。

**请求URL**: `/health`

**Host**: `10.5.122.114:8080`

**请求方式**: `GET`

**headers**: 无

**uri参数**: 无

**body参数**: 无

**返回参数**: 参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| status | string | 服务状态 |
| message | string | 状态描述信息 |

**错误码**: 无

**调用示例**:

**请求示例**:
```bash
curl http://10.5.122.114:8080/health
```

**返回示例**:
```json
{
  "status": "ok",
  "message": "服务正常运行"
}
```

---

## 2. 配置查看接口

**简要描述**: 查看当前预警系统的配置信息，包括邮件配置等。

**请求URL**: `/config`

**Host**: `10.5.122.114:8080`

**请求方式**: `GET`

**headers**: 无

**uri参数**: 无

**body参数**: 无

**返回参数**: 参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| status | string | 响应状态 |
| email_config | object | 邮件配置信息 |
| email_config.api_url | string | 邮件API地址 |
| email_config.app_id | string | 应用ID |
| email_config.app_secret | string | 应用密钥（隐藏） |
| email_config.from | string | 发件人地址 |
| email_config.to | array | 收件人地址列表 |
| email_config.debug_mode | boolean | 调试模式 |
| email_config.debug_api_url | string | 调试API地址 |

**错误码**: 无

**调用示例**:

**请求示例**:
```bash
curl http://10.5.122.114:8080/config
```

**返回示例**:
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

---

## 3. 邮件测试接口

**简要描述**: 测试预警系统的邮件发送功能，发送包含测试预警信息的邮件。

**请求URL**: `/test-email`

**Host**: `10.5.122.114:8080`

**请求方式**: `POST`

**headers**: 无

**uri参数**: 无

**body参数**: 无

**返回参数**: 参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | integer | 响应状态码 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.alert_count | integer | 测试预警数量 |
| data.recipients | array | 收件人列表 |

**错误码**:

| 错误码 | 说明 |
|--------|------|
| 500 | 邮件发送失败 |

**调用示例**:

**请求示例**:
```bash
curl -X POST http://10.5.122.114:8080/test-email
```

**返回示例**:
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

---

## 4. 创建预警信息接口

**简要描述**: 向预警系统提交新的预警信息，系统将自动存储于MySQL数据库并可用于后续的邮件通知。

**请求URL**: `/api/v1/alerts`

**Host**: `10.5.122.114:8080`

**请求方式**: `POST`

**headers**: 

| 参数名 | 必选 | 说明 |
|--------|------|------|
| Content-Type | 是 | 请求体格式，固定值：application/json |

**uri参数**: 无

**body参数[json]**: 

| 参数名 | 必选 | 类型 | 说明 |
|--------|------|------|------|
| domain | 是 | string | 域名 |
| message | 是 | string | 预警信息内容 |
| source | 是 | string | 预警来源 |
| status | 否 | string | 预警状态，如 "active", "resolved", "warning" |
| region | 否 | string | 区域信息 |
| alert_time | 否 | string | 预警时间，格式为 "YYYY-MM-DD HH:mm:ss"，默认为当前时间 |

**返回参数**: 参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | integer | 响应状态码 |
| message | string | 响应消息 |
| data | object | 预警信息数据 |
| data.id | integer | 预警ID |
| data.domain | string | 域名 |
| data.message | string | 预警信息 |
| data.source | string | 预警来源 |
| data.status | string | 预警状态 |
| data.region | string | 区域信息 |
| data.alert_time | string | 预警时间 |
| data.created_at | string | 创建时间 |
| data.updated_at | string | 更新时间 |

**错误码**:

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 500 | 存储预警信息失败 |

**调用示例**:

**请求示例**:
```bash
curl -X POST "http://10.5.122.114:8080/api/v1/alerts" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "search.suggest.kgidc.cn",
    "message": "北方已切量,但南方超过24小时未切量,请检查",
    "source": "RPC后台",
    "status": "active",
    "region": "全国"
  }'
```

**返回示例**:
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

---

## 5. 获取预警信息接口

**简要描述**: 获取系统中的所有预警信息，支持分页查询。

**请求URL**: `/api/v1/alerts`

**Host**: `10.5.122.114:8080`

**请求方式**: `GET`

**headers**: 无

**uri参数**: 

| 参数名 | 必选 | 类型 | 说明 |
|--------|------|------|------|
| page | 否 | integer | 页码，默认为1 |
| size | 否 | integer | 每页数量，默认为20，最大100 |

**body参数**: 无

**返回参数**: 参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | integer | 响应状态码 |
| message | string | 响应消息 |
| data | array | 预警信息列表 |
| data[].id | integer | 预警ID |
| data[].domain | string | 域名 |
| data[].message | string | 预警信息 |
| data[].source | string | 预警来源 |
| data[].status | string | 预警状态 |
| data[].region | string | 区域信息 |
| data[].alert_time | string | 预警时间 |
| data[].created_at | string | 创建时间 |
| data[].updated_at | string | 更新时间 |
| total | integer | 总记录数 |
| page | integer | 当前页码 |
| size | integer | 每页大小 |

**错误码**:

| 错误码 | 说明 |
|--------|------|
| 500 | 获取预警信息失败 |

**调用示例**:

**请求示例**:
```bash
curl "http://10.5.122.114:8080/api/v1/alerts?page=1&size=20"
```

**返回示例**:
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

---

## 6. 按时间段查询预警接口

**简要描述**: 根据指定的时间范围查询预警信息，支持自定义开始和结束时间。

**请求URL**: `/api/v1/alerts/period`

**Host**: `10.5.122.114:8080`

**请求方式**: `GET`

**headers**: 无

**uri参数**: 

| 参数名 | 必选 | 类型 | 说明 |
|--------|------|------|------|
| start_time | 否 | string | 开始时间，格式为 "YYYY-MM-DD HH:mm:ss"，默认为当天晚上7点 |
| end_time | 否 | string | 结束时间，格式为 "YYYY-MM-DD HH:mm:ss"，默认为当天晚上10点 |

**body参数**: 无

**返回参数**: 参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | integer | 响应状态码 |
| message | string | 响应消息 |
| data | array | 预警信息列表 |
| data[].id | integer | 预警ID |
| data[].domain | string | 域名 |
| data[].message | string | 预警信息 |
| data[].source | string | 预警来源 |
| data[].status | string | 预警状态 |
| data[].region | string | 区域信息 |
| data[].alert_time | string | 预警时间 |
| data[].created_at | string | 创建时间 |
| data[].updated_at | string | 更新时间 |
| start_time | string | 查询开始时间 |
| end_time | string | 查询结束时间 |
| total | integer | 总记录数 |

**错误码**:

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 500 | 获取预警信息失败 |

**调用示例**:

**请求示例**:
```bash
curl "http://10.5.122.114:8080/api/v1/alerts/period?start_time=2025-01-15%2007:00:00&end_time=2025-01-15%2022:00:00"
```

**返回示例**:
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

---

## 通用说明

### 系统信息
- **服务名称**: 预警系统 API
- **部署服务器**: 10.5.122.114 (CentOS 7)
- **数据库**: 广州MySQL (10.5.122.136:3306)
- **邮件服务**: HTTP API方式

### 注意事项
1. 所有时间参数格式必须为 "YYYY-MM-DD HH:mm:ss"
2. 如果不提供预警时间，系统将自动使用当前时间
3. 系统会在每天晚上10点自动统计并发送邮件通知
4. 所有接口返回的JSON数据均使用UTF-8编码 