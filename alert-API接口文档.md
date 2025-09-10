预警系统 API 接口文档

## 1. 健康检查接口

一、简要描述
检查预警系统服务是否正常运行,用于监控和健康检查。

二、请求URL
http://10.5.122.114:8080/health

三、Host
预发布环境: 10.5.122.114:8080

四、请求方式
GET

五、headers
无

六、uri参数
无

七、body参数
无

八、返回参数
参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| status | string | 服务状态 |
| message | string | 状态描述信息 |

九、错误码
无

十、调用示例

请求示例:
```bash
curl http://10.5.122.114:8080/health
```

返回示例:
```json
{
  "status": "ok",
  "message": "服务正常运行"
}
```

---

## 2. 配置查看接口

一、简要描述
查看当前预警系统的配置信息，包括邮件配置和定时任务配置。

二、请求URL
http://10.5.122.114:8080/config

三、Host
预发布环境: 10.5.122.114:8080

四、请求方式
GET

五、headers
无

六、uri参数
无

七、body参数
无

八、返回参数
参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| status | string | 响应状态 |
| email_config | object | 邮件配置信息 |
| email_config.api_url | string | 邮件API地址 |
| email_config.app_id | string | 应用ID |
| email_config.app_secret | string | 应用密钥（隐藏） |
| email_config.from | string | 发件人地址 |
| email_config.debug_mode | boolean | 调试模式 |
| email_config.debug_api_url | string | 调试API地址 |
| email_config.note | string | 说明信息 |
| cron_config | object | 定时任务配置信息 |
| cron_config.enabled | boolean | 是否启用定时任务 |
| cron_config.schedule | string | Cron表达式 |
| cron_config.start_time | string | 查询开始时间 |
| cron_config.end_time | string | 查询结束时间 |
| cron_config.description | string | 配置说明 |

九、错误码
无

十、调用示例

请求示例:
```bash
curl http://10.5.122.114:8080/config
```

返回示例:
```json
{
  "status": "ok",
  "email_config": {
    "api_url": "http://opi.kgidc.cn/mail/email/send_email.php",
    "app_id": "v1-5f4769fe10c9c",
    "app_secret": "***hidden***",
    "from": "system@company.com",
    "debug_mode": false,
    "debug_api_url": "http://10.16.2.146:6709/mail/email/send_email.php",
    "note": "收件人现在根据告警信息动态生成"
  },
  "cron_config": {
    "enabled": true,
    "schedule": "0 22 * * *",
    "start_time": "19:00",
    "end_time": "22:00",
    "description": "定时任务配置信息"
  }
}
```

---

## 3. 邮件测试接口

一、简要描述
测试预警系统的邮件发送功能，发送包含测试预警信息的邮件。系统会自动生成测试数据并发送给多个测试用户。

二、请求URL
http://10.5.122.114:8080/test-email

三、Host
预发布环境: 10.5.122.114:8080

四、请求方式
POST

五、headers
无

六、uri参数
无

七、body参数
无

八、返回参数
参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | integer | 响应状态码 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.user_count | integer | 用户数量 |
| data.total_alerts | integer | 测试预警数量 |
| data.recipients | array | 收件人列表 |
| data.details | array | 详细发送信息 |
| data.details[].user | string | 用户名 |
| data.details[].alert_count | integer | 该用户的预警数量 |

九、错误码

| 错误码 | 说明 |
|--------|------|
| 500 | 邮件发送失败 |

十、调用示例

请求示例:
```bash
curl -X POST http://10.5.122.114:8080/test-email
```

返回示例:
```json
{
  "code": 200,
  "message": "测试邮件发送成功",
  "data": {
    "user_count": 4,
    "total_alerts": 7,
    "recipients": ["felixgao", "hugoli", "zhangsan", "lisi"],
    "details": [
      {"user": "felixgao", "alert_count": 2},
      {"user": "hugoli", "alert_count": 2},
      {"user": "zhangsan", "alert_count": 1},
      {"user": "lisi", "alert_count": 2}
    ]
  }
}
```

---
## 4. 创建预警信息接口

一、简要描述
向预警系统提交新的预警信息，系统将自动存储并可用于后续的邮件通知。收件人信息将用于动态生成邮件地址。

二、请求URL
http://10.5.122.114:8080/api/v1/alerts

三、Host
预发布环境: 10.5.122.114:8080

四、请求方式
POST

五、headers

| 参数名 | 必选 | 说明 |
|--------|------|------|
| Content-Type | 是 | 请求体格式，固定值：application/json |

六、uri参数
无

七、body参数[json]

| 参数名 | 必选 | 类型 | 说明 |
|--------|------|------|------|
| message | 是 | string | 预警信息内容 |
| recipient | 是 | string | 收件人标识，系统会自动添加@kugou.net后缀生成邮箱地址 |
| alert_time | 否 | string | 预警时间，格式为 "YYYY-MM-DD HH:mm:ss"，默认为当前时间 |

八、返回参数
参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | integer | 响应状态码 |
| message | string | 响应消息 |
| data | object | 预警信息数据 |
| data.id | integer | 预警ID |
| data.message | string | 预警信息 |
| data.recipient | string | 收件人标识 |
| data.alert_time | string | 预警时间 |
| data.created_at | string | 创建时间 |
| data.updated_at | string | 更新时间 |

九、错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 500 | 存储预警信息失败 |

十、调用示例

请求示例:
```bash
curl -X POST "http://10.5.122.114:8080/api/v1/alerts" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
    "recipient": "zhangsan"
  }'
```

返回示例:
```json
{
  "code": 200,
  "message": "预警信息创建成功",
  "data": {
    "id": 1,
    "message": "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
    "recipient": "zhangsan",
    "alert_time": "2025-01-15T19:30:00+08:00",
    "created_at": "2025-01-15T19:30:00+08:00",
    "updated_at": "2025-01-15T19:30:00+08:00"
  }
}
```

---

## 5. 获取预警信息接口

一、简要描述
获取系统中的所有预警信息，支持分页查询。

二、请求URL
http://10.5.122.114:8080/api/v1/alerts

三、Host
预发布环境: 10.5.122.114:8080

四、请求方式
GET

五、headers
无

六、uri参数

| 参数名 | 必选 | 类型 | 说明 |
|--------|------|------|------|
| page | 否 | integer | 页码，默认为1 |
| page_size | 否 | integer | 每页数量，默认为20，最大100 |

七、body参数
无

八、返回参数
参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | integer | 响应状态码 |
| message | string | 响应消息 |
| data | array | 预警信息列表 |
| data[].id | integer | 预警ID |
| data[].message | string | 预警信息 |
| data[].recipient | string | 收件人标识 |
| data[].alert_time | string | 预警时间 |
| data[].created_at | string | 创建时间 |
| data[].updated_at | string | 更新时间 |
| total | integer | 总记录数 |
| page | integer | 当前页码 |
| size | integer | 每页大小 |

九、错误码

| 错误码 | 说明 |
|--------|------|
| 500 | 获取预警信息失败 |

十、调用示例

请求示例:
```bash
curl "http://10.5.122.114:8080/api/v1/alerts?page=1&page_size=20"
```

返回示例:
```json
{
  "code": 200,
  "message": "获取预警信息成功",
  "data": [
    {
      "id": 1,
      "message": "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
      "recipient": "zhangsan",
      "alert_time": "2025-01-15T19:30:00+08:00",
      "created_at": "2025-01-15T19:30:00+08:00",
      "updated_at": "2025-01-15T19:30:00+08:00"
    }
  ],
  "total": 1,
  "page": 1,
  "size": 20
}
```

---

## 6. 按收件人查询预警接口

一、简要描述
根据指定的收件人查询预警信息，支持按用户分组查看。

二、请求URL
http://10.5.122.114:8080/api/v1/alerts/recipient

三、Host
预发布环境: 10.5.122.114:8080

四、请求方式
GET

五、headers
无

六、uri参数

| 参数名 | 必选 | 类型 | 说明 |
|--------|------|------|------|
| recipient | 是 | string | 收件人标识 |

七、body参数
无

八、返回参数
参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | integer | 响应状态码 |
| message | string | 响应消息 |
| data | array | 预警信息列表 |
| data[].id | integer | 预警ID |
| data[].message | string | 预警信息 |
| data[].recipient | string | 收件人标识 |
| data[].alert_time | string | 预警时间 |
| data[].created_at | string | 创建时间 |
| data[].updated_at | string | 更新时间 |
| recipient | string | 查询的收件人 |
| total | integer | 总记录数 |

九、错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 收件人参数不能为空 |
| 500 | 获取预警信息失败 |

十、调用示例

请求示例:
```bash
curl "http://10.5.122.114:8080/api/v1/alerts/recipient?recipient=zhangsan"
```

返回示例:
```json
{
  "code": 200,
  "message": "获取预警信息成功",
  "data": [
    {
      "id": 1,
      "message": "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
      "recipient": "zhangsan",
      "alert_time": "2025-01-15T19:30:00+08:00",
      "created_at": "2025-01-15T19:30:00+08:00",
      "updated_at": "2025-01-15T19:30:00+08:00"
    }
  ],
  "recipient": "zhangsan",
  "total": 1
}
```

---

## 7. 按时间段查询预警接口

一、简要描述
根据指定的时间范围查询预警信息，支持自定义开始和结束时间。

二、请求URL
http://10.5.122.114:8080/api/v1/alerts/period

三、Host
预发布环境: 10.5.122.114:8080

四、请求方式
GET

五、headers
无

六、uri参数

| 参数名 | 必选 | 类型 | 说明 |
|--------|------|------|------|
| start_time | 否 | string | 开始时间，格式为 "YYYY-MM-DD HH:mm:ss"，默认为当天晚上7点 |
| end_time | 否 | string | 结束时间，格式为 "YYYY-MM-DD HH:mm:ss"，默认为当天晚上10点 |

七、body参数
无

八、返回参数
参数以json形式返回

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | integer | 响应状态码 |
| message | string | 响应消息 |
| data | array | 预警信息列表 |
| data[].id | integer | 预警ID |
| data[].message | string | 预警信息 |
| data[].recipient | string | 收件人标识 |
| data[].alert_time | string | 预警时间 |
| data[].created_at | string | 创建时间 |
| data[].updated_at | string | 更新时间 |
| start_time | string | 查询开始时间 |
| end_time | string | 查询结束时间 |
| total | integer | 总记录数 |

九、错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 500 | 获取预警信息失败 |

十、调用示例

请求示例:
```bash
curl "http://10.5.122.114:8080/api/v1/alerts/period?start_time=2025-01-15%2019:00:00&end_time=2025-01-15%2022:00:00"
```

返回示例:
```json
{
  "code": 200,
  "message": "获取预警信息成功",
  "data": [
    {
      "id": 1,
      "message": "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
      "recipient": "zhangsan",
      "alert_time": "2025-01-15T19:30:00+08:00",
      "created_at": "2025-01-15T19:30:00+08:00",
      "updated_at": "2025-01-15T19:30:00+08:00"
    }
  ],
  "start_time": "2025-01-15 19:00:00",
  "end_time": "2025-01-15 22:00:00",
  "total": 1
}
```

---

## 通用说明

### 系统信息
- 服务名称: 预警系统 API
- 部署服务器: 10.5.122.114 (CentOS 7)
- 数据库: 广州MySQL (10.5.122.136:3306) - alert_message
- 邮件服务: HTTP API方式

### 核心功能特性
1. **动态收件人**: 系统根据告警信息中的recipient字段自动生成邮箱地址（添加@kugou.net后缀）
2. **用户分组**: 定时任务按收件人分组发送邮件，每个用户收到专属的告警信息
3. **简化结构**: 只使用message和recipient两个核心字段，简化了数据结构
4. **自动邮件**: 每天晚上10点自动统计当天晚上7点到10点的告警信息并发送邮件
5. **用户列表管理**: 支持从userlist.json文件加载用户信息，自动映射英文名到邮箱地址
6. **管理员邮件**: 当用户未找到时，自动发送合并邮件给管理员（liyongchang@kugou.net）

### 邮件功能说明
1. **用户邮件**: 为每个用户生成专属的HTML邮件，包含其所有预警信息
2. **管理员邮件**: 当用户未找到时，发送包含所有未找到用户预警信息的合并邮件给liyongchang@kugou.net
3. **邮件模板**: 美观的HTML格式，包含预警概览、详细信息和时间范围
4. **编码支持**: 完整支持UTF-8编码，确保中文内容正确显示

### 注意事项
1. 所有时间参数格式必须为 "YYYY-MM-DD HH:mm:ss"
2. 如果不提供预警时间，系统将自动使用当前时间
3. 收件人字段会自动添加@kugou.net后缀生成邮箱地址
4. 系统会在每天晚上10点自动统计并发送邮件通知
5. 所有接口返回的JSON数据均使用UTF-8编码
6. 邮件发送按用户分组，每个用户只收到属于自己的告警信息
7. 邮件中不显示预警ID，只显示预警编号、内容和时间
8. 支持调试模式，可以配置不同的邮件API地址进行测试