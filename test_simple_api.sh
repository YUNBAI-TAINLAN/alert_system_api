#!/bin/bash

# 告警系统API测试脚本（简化版，不使用jq）
# 使用方法: ./test_simple_api.sh [base_url]
# 默认base_url: http://localhost:8080

BASE_URL=${1:-"http://localhost:8080"}

echo "🚀 开始测试告警系统API（简化版）..."
echo "📍 目标地址: $BASE_URL"
echo ""

# 测试健康检查
echo "1️⃣ 测试健康检查接口..."
curl -s "$BASE_URL/health"
echo ""
echo ""

# 测试创建告警信息
echo "2️⃣ 测试创建告警信息..."
curl -s -X POST "$BASE_URL/api/v1/alerts" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.com",
    "message": "服务器CPU使用率超过90%",
    "source": "监控系统",
    "status": "critical",
    "region": "北京"
  }'
echo ""
echo ""

# 再创建几条测试数据
echo "3️⃣ 创建更多测试数据..."
for i in {1..3}; do
curl -s -X POST "$BASE_URL/api/v1/alerts" \
  -H "Content-Type: application/json" \
  -d '{
      "domain": "test'$i'.com",
      "message": "测试告警信息 #'$i'",
      "source": "测试系统",
      "status": "warning",
      "region": "上海"
    }' > /dev/null
done
echo "✅ 创建了3条测试数据"
echo ""

# 测试获取告警信息
echo "4️⃣ 测试获取告警信息..."
curl -s "$BASE_URL/api/v1/alerts?page=1&page_size=5"
echo ""
echo ""

# 测试按时间段查询
echo "5️⃣ 测试按时间段查询告警..."
START_TIME=$(date -d "1 hour ago" "+%Y-%m-%d %H:%M:%S")
END_TIME=$(date "+%Y-%m-%d %H:%M:%S")
curl -s "$BASE_URL/api/v1/alerts/period?start_time=$START_TIME&end_time=$END_TIME"
echo ""
echo ""

echo " API测试完成！"
echo ""
echo " 提示："
echo "   - 如果看到JSON响应，说明API工作正常"
echo "   - 如果看到错误，请确保服务正在运行"
echo "   - 检查数据库连接配置"
