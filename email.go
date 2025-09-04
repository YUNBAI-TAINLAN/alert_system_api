package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// EmailConfig 邮件配置
type EmailConfig struct {
	APIUrl      string   `json:"api_url"`
	AppID       string   `json:"app_id"`
	AppSecret   string   `json:"app_secret"`
	From        string   `json:"from"`
	To          []string `json:"to"`
	DebugMode   bool     `json:"debug_mode"`
	DebugAPIUrl string   `json:"debug_api_url"`
}

// EmailAPIRequest 邮件API请求结构
type EmailAPIRequest struct {
	OpdAppid    string `json:"opdAppid"`
	OpdAppsecret string `json:"opdAppsecret"`
	ToList      string `json:"to_list"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	Mimetype    string `json:"mimetype"`
}

// EmailAPIResponse 邮件API响应结构
type EmailAPIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

var emailConfig EmailConfig

// InitEmailConfig 初始化邮件配置
func InitEmailConfig() {
	// 使用全局配置
	emailConfig = config.Email
}

// SendAlertEmail 发送预警通知邮件（按用户分组）
func SendAlertEmail(userAlertsList []UserAlerts) error {
	if len(userAlertsList) == 0 {
		return fmt.Errorf("没有预警信息需要发送")
	}

	log.Printf("开始发送邮件，共涉及 %d 个用户", len(userAlertsList))
	
	var successCount, failCount int
	var successRecipients, failRecipients []string

	// 为每个用户发送单独的邮件
	for _, userAlerts := range userAlertsList {
		recipientEmail := generateRecipientEmail(userAlerts.Recipient)
		log.Printf("正在发送邮件给用户: %s (%s)，包含 %d 条预警信息", 
			userAlerts.Recipient, recipientEmail, len(userAlerts.Alerts))
		
		if err := sendEmailToUser(userAlerts); err != nil {
			log.Printf("❌ 发送邮件给用户 %s (%s) 失败: %v", 
				userAlerts.Recipient, recipientEmail, err)
			failCount++
			failRecipients = append(failRecipients, recipientEmail)
			continue
		}
		
		log.Printf("✅ 成功发送邮件给用户: %s (%s)，包含 %d 条预警信息", 
			userAlerts.Recipient, recipientEmail, len(userAlerts.Alerts))
		successCount++
		successRecipients = append(successRecipients, recipientEmail)
	}

	// 发送总结
	log.Printf("📧 邮件发送完成:")
	log.Printf("   ✅ 成功: %d 个用户", successCount)
	log.Printf("   ❌ 失败: %d 个用户", failCount)
	if len(successRecipients) > 0 {
		log.Printf("   📬 成功收件人: %v", successRecipients)
	}
	if len(failRecipients) > 0 {
		log.Printf("   📭 失败收件人: %v", failRecipients)
	}

	if failCount > 0 {
		return fmt.Errorf("部分邮件发送失败，成功: %d，失败: %d", successCount, failCount)
	}

	return nil
}

// sendEmailToUser 发送邮件给特定用户
func sendEmailToUser(userAlerts UserAlerts) error {
	// 生成邮件内容
	subject, body, err := generateEmailContentForUser(userAlerts)
	if err != nil {
		return fmt.Errorf("生成邮件内容失败: %v", err)
	}

	// 生成收件人邮箱地址
	recipientEmail := generateRecipientEmail(userAlerts.Recipient)

	// 通过HTTP API发送邮件
	if err := sendEmailViaAPI([]string{recipientEmail}, subject, body); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

// generateRecipientEmail 根据收件人变量生成邮箱地址
func generateRecipientEmail(recipient string) string {
	// 如果收件人已经包含@符号，直接返回
	if strings.Contains(recipient, "@") {
		return recipient
	}
	
	// 否则添加@kugou.net后缀
	return recipient + "@kugou.net"
}

// sendEmailViaAPI 通过HTTP API发送邮件
func sendEmailViaAPI(toUsers []string, subject, content string) error {
	// 确定使用的API地址
	apiURL := emailConfig.APIUrl
	if emailConfig.DebugMode && emailConfig.DebugAPIUrl != "" {
		apiURL = emailConfig.DebugAPIUrl
		log.Printf("🔧 使用调试模式邮件API: %s", apiURL)
	}

		// 将收件人列表转换为逗号分隔的字符串
	toList := strings.Join(toUsers, ",")

	// 构建表单数据（与您同事的PHP代码保持一致）
	formData := url.Values{}
	formData.Set("opdAppid", emailConfig.AppID)
	formData.Set("opdAppsecret", emailConfig.AppSecret)
	formData.Set("to_list", toList)
	formData.Set("subject", subject)
	formData.Set("body", content)
	formData.Set("mimetype", "html")

	// 将表单数据转换为字符串
	postData := formData.Encode()

	// 创建HTTP请求
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(postData))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头（使用表单数据格式）
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Alert-System/1.0")

	// 记录请求详情（简化版，避免敏感信息泄露）
	log.Printf("📤 发送邮件请求:")
	log.Printf("    收件人: %s", toList)
	log.Printf("    主题: %s", subject)
	log.Printf("    API地址: %s", apiURL)

	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	// 记录响应详情
	log.Printf("📥 邮件API响应:")
	log.Printf("   📊 状态码: %s", resp.Status)
	log.Printf("   📄 响应内容: %s", string(respBody))

	// 解析响应
	var apiResp EmailAPIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		// 如果解析JSON失败，记录原始响应
		log.Printf("⚠️ 邮件API响应解析失败，原始响应: %s", string(respBody))
		return fmt.Errorf("解析API响应失败: %v", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("邮件API返回错误状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 检查API响应中的业务状态码
	if apiResp.Code != 0 && apiResp.Code != 200 { // 假设0或200表示成功
		return fmt.Errorf("邮件发送失败: code=%d, message=%s", apiResp.Code, apiResp.Message)
	}

	log.Printf("✅ 邮件发送成功: %s", apiResp.Message)
	return nil
}

// generateEmailContentForUser 为用户生成邮件内容
func generateEmailContentForUser(userAlerts UserAlerts) (string, string, error) {
	// 邮件主题
	subject := fmt.Sprintf("预警通知 - %s - %s", userAlerts.Recipient, time.Now().Format("2006-01-02"))

	// 邮件模板 - 预警通知样式
	const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>预警通知</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; 
            margin: 0; 
            padding: 20px; 
            background-color: #f5f5f5; 
            line-height: 1.6;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
            font-weight: 600;
        }
        .header p {
            margin: 10px 0 0 0;
            opacity: 0.9;
            font-size: 14px;
        }
        .content {
            padding: 30px;
        }
        .summary {
            background-color: #f8f9fa;
            border-left: 4px solid #007bff;
            padding: 20px;
            margin-bottom: 30px;
            border-radius: 0 8px 8px 0;
        }
        .summary h3 {
            margin: 0 0 15px 0;
            color: #333;
            font-size: 18px;
        }
        .summary ul {
            margin: 0;
            padding-left: 20px;
        }
        .summary li {
            margin-bottom: 8px;
            color: #555;
        }
        .alert-card {
            background: #ffffff;
            border: 1px solid #e1e5e9;
            border-radius: 12px;
            margin-bottom: 20px;
            padding: 20px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.06);
            transition: transform 0.2s ease, box-shadow 0.2s ease;
        }
        .alert-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 16px rgba(0,0,0,0.12);
        }
        .alert-header {
            display: flex;
            align-items: center;
            margin-bottom: 15px;
        }
        .alert-icon {
            width: 40px;
            height: 40px;
            background: linear-gradient(135deg, #28a745, #20c997);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 15px;
            flex-shrink: 0;
        }
        .alert-icon::before {
            content: "!";
            color: white;
            font-size: 18px;
            font-weight: bold;
        }
        .alert-title {
            flex: 1;
        }
        .alert-title h4 {
            margin: 0;
            color: #333;
            font-size: 16px;
            font-weight: 600;
        }
        .alert-time {
            color: #666;
            font-size: 13px;
            margin-top: 2px;
        }
        .alert-content {
            margin-left: 55px;
        }
        .alert-message {
            background-color: #fff3cd;
            border: 1px solid #ffeaa7;
            border-radius: 6px;
            padding: 12px;
            margin-bottom: 10px;
            line-height: 1.5;
        }
        .footer {
            background-color: #f8f9fa;
            padding: 20px 30px;
            text-align: center;
            border-top: 1px solid #e1e5e9;
            color: #666;
            font-size: 13px;
        }
    </style>
</head>
<body>
    <div class="container">
            <div class="header">
            <h1>预警通知</h1>
            <p>收件人: {{.Recipient}} | 生成时间: {{.GenerateTime}}</p>
        </div>
    
        <div class="content">
    <div class="summary">
                    <h3>统计摘要</h3>
        <ul>
                    <li>总预警数量: <strong>{{.TotalCount}}</strong></li>
                    <li>收件人: <strong>{{.Recipient}}</strong></li>
                    <li>统计时间段: {{.StartTime}} 至 {{.EndTime}}</li>
        </ul>
    </div>

            <h3 style="margin-bottom: 20px; color: #333; font-size: 18px;">详细预警信息</h3>
            
            {{range $index, $alert := .Alerts}}
            <div class="alert-card">
                <div class="alert-header">
                    <div class="alert-icon"></div>
                    <div class="alert-title">
                        <h4>预警通知 {{$alert.AlertTime.Format "1/2 15:04:05"}}</h4>
                        <div class="alert-time">{{$alert.AlertTime.Format "2006-01-02 15:04:05"}}</div>
                    </div>
                </div>
                <div class="alert-content">
                    <div class="alert-message">
                        {{$alert.Message}}
                    </div>
                        <div class="detail-item">
                            <span class="detail-label">时间:</span>
                            <span>{{$alert.AlertTime.Format "2006-01-02 15:04:05"}}</span>
                    </div>
                </div>
            </div>
            {{end}}
        </div>
    
        <div class="footer">
            <p><strong>注意:</strong> 此邮件为系统自动生成，请及时处理相关预警信息。</p>
        <p>如有疑问，请联系系统管理员。</p>
        </div>
    </div>
</body>
</html>
`

	// 准备模板数据
	type TemplateData struct {
		GenerateTime string
		StartTime    string
		EndTime      string
		TotalCount   int
		Recipient    string
		Alerts       []Alert
	}

	// 确定时间范围
	var startTime, endTime time.Time
	if len(userAlerts.Alerts) > 0 {
		startTime = userAlerts.Alerts[len(userAlerts.Alerts)-1].AlertTime
		endTime = userAlerts.Alerts[0].AlertTime
	}

	data := TemplateData{
		GenerateTime: time.Now().Format("2006-01-02 15:04:05"),
		StartTime:    startTime.Format("2006-01-02 15:04:05"),
		EndTime:      endTime.Format("2006-01-02 15:04:05"),
		TotalCount:   len(userAlerts.Alerts),
		Recipient:    userAlerts.Recipient,
		Alerts:       userAlerts.Alerts,
	}

	// 解析模板
	tmpl, err := template.New("email").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}).Parse(emailTemplate)
	if err != nil {
		return "", "", fmt.Errorf("解析邮件模板失败: %v", err)
	}

	// 执行模板
	var body strings.Builder
	if err := tmpl.Execute(&body, data); err != nil {
		return "", "", fmt.Errorf("执行邮件模板失败: %v", err)
	}

	return subject, body.String(), nil
} 