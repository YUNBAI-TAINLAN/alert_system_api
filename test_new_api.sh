#!/bin/bash

# 告警系统API测试脚本（新版本）
# 使用方法: ./test_new_api.sh [base_url]
# 默认base_url: http://localhost:8080

BASE_URL=${1:-"http://localhost:8080"}

echo "🚀 开始测试告警系统API（新版本）..."
echo "📍 目标地址: $BASE_URL"
echo ""

# 测试健康检查
echo "1️⃣ 测试健康检查接口..."
curl -s "$BASE_URL/health"
echo ""
echo ""

# 测试创建告警信息（新格式）- 单个收件人
echo "2️⃣ 测试创建告警信息（单个收件人）..."
curl -s -X POST "$BASE_URL/api/v1/alerts" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "检测到域名【search.suggest.kgidc.cn】北方已切量，但南方超过24小时未切量，请检查",
    "recipient": "felixgao"
  }'
echo ""
echo ""

# 测试创建告警信息（新格式）- 多个收件人（英文逗号）
echo "3️⃣ 测试创建告警信息（多个收件人-英文逗号）..."
curl -s -X POST "$BASE_URL/api/v1/alerts" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "检测到域名【api.example.com】响应时间超过5秒",
    "recipient": "zhangsan,lisi"
  }'
echo ""
echo ""

# 测试创建告警信息（新格式）- 多个收件人（中文逗号）
echo "4️⃣ 测试创建告警信息（多个收件人-中文逗号）..."
curl -s -X POST "$BASE_URL/api/v1/alerts" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "检测到域名【web.example.com】连接数达到上限",
    "recipient": "wangwu，maliu"
  }'
echo ""
echo ""

# 测试创建告警信息（新格式）- 多个收件人（混合逗号）
echo "5️⃣ 测试创建告警信息（多个收件人-混合逗号）..."
curl -s -X POST "$BASE_URL/api/v1/alerts" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "检测到域名【cdn.kugou.com】CDN节点异常",
    "recipient": "qianqi，sunba,zhoujiu"
  }'
echo ""
echo ""

echo "✅ 创建了多条测试数据，包含单个和多个收件人"
echo ""

# 测试获取告警信息
echo "6️⃣ 测试获取告警信息..."
curl -s "$BASE_URL/api/v1/alerts?page=1&page_size=10"
echo ""
echo ""

# 测试按收件人查询
echo "7️⃣ 测试按收件人查询告警..."
curl -s "$BASE_URL/api/v1/alerts/recipient?recipient=felixgao"
echo ""
echo ""

# 测试按收件人查询（多个收件人中的一个）
echo "8️⃣ 测试按收件人查询告警（zhangsan）..."
curl -s "$BASE_URL/api/v1/alerts/recipient?recipient=zhangsan"
echo ""
echo ""

# 测试按时间段查询
echo "9️⃣ 测试按时间段查询告警..."
START_TIME=$(date -d "1 hour ago" "+%Y-%m-%d %H:%M:%S")
END_TIME=$(date "+%Y-%m-%d %H:%M:%S")
curl -s "$BASE_URL/api/v1/alerts/period?start_time=$START_TIME&end_time=$END_TIME"
echo ""
echo ""

# 测试邮件发送功能
echo "🔟 测试邮件发送功能..."
curl -s -X POST "$BASE_URL/test-email" \
  -H "Content-Type: application/json"
echo ""
echo ""

echo "✅ API测试完成！"
echo ""
echo "�� 测试总结："
echo "   - 新API只使用message和recipient字段"
echo "   - 支持单个收件人和多个收件人（逗号分隔）"
echo "   - 支持中文逗号和英文逗号分隔"
echo "   - 支持按收件人查询告警信息"
echo "   - 邮件发送按用户分组，动态生成收件人邮箱"
echo "   - 定时任务按用户分组发送告警信息"
echo ""
echo "�� 提示："
echo "   - 如果看到JSON响应，说明API工作正常"
echo "   - 如果看到错误，请确保服务正在运行"
echo "   - 检查数据库连接配置"
echo "   - 检查邮件配置是否正确"
echo "   - 多个收件人会自动创建多条记录，每个收件人一条"