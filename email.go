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

// SendAlertEmail 发送预警通知邮件
func SendAlertEmail(alerts []Alert) error {
	if len(alerts) == 0 {
		return fmt.Errorf("没有预警信息需要发送")
	}

	// 生成邮件内容
	subject, body, err := generateEmailContent(alerts)
	if err != nil {
		return fmt.Errorf("生成邮件内容失败: %v", err)
	}

	// 通过HTTP API发送邮件
	if err := sendEmailViaAPI(emailConfig.To, subject, body); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	log.Printf("成功发送预警通知邮件给 %d 个收件人", len(emailConfig.To))
	return nil
}

// sendEmailViaAPI 通过HTTP API发送邮件
func sendEmailViaAPI(toUsers []string, subject, content string) error {
	// 确定使用的API地址
	apiURL := emailConfig.APIUrl
	if emailConfig.DebugMode && emailConfig.DebugAPIUrl != "" {
		apiURL = emailConfig.DebugAPIUrl
		log.Printf("使用调试模式邮件API: %s", apiURL)
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

	// 记录请求详情
	log.Printf("发送邮件API请求:")
	log.Printf("  URL: %s", apiURL)
	log.Printf("  Method: %s", req.Method)
	log.Printf("  Headers: %v", req.Header)
	log.Printf("  Body: %s", postData)

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
	log.Printf("邮件API响应:")
	log.Printf("  Status: %s", resp.Status)
	log.Printf("  Headers: %v", resp.Header)
	log.Printf("  Body: %s", string(respBody))

	// 解析响应
	var apiResp EmailAPIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		// 如果解析JSON失败，记录原始响应
		log.Printf("邮件API响应解析失败，原始响应: %s", string(respBody))
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

	log.Printf("邮件发送成功: %s", apiResp.Message)
	return nil
}

// generateEmailContent 生成邮件内容
func generateEmailContent(alerts []Alert) (string, string, error) {
	// 邮件主题
	subject := fmt.Sprintf("预警通知汇总 - %s", time.Now().Format("2006-01-02"))

	// 邮件模板 - 预警通知样式
	const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>预警通知汇总</title>
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
        .domain-highlight {
            background-color: #fff3cd;
            color: #0056b3;
            font-weight: 600;
            padding: 2px 6px;
            border-radius: 4px;
            margin: 0 4px;
        }
        .alert-details {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 10px;
            font-size: 13px;
            color: #666;
        }
        .detail-item {
            display: flex;
            align-items: center;
        }
        .detail-label {
            font-weight: 600;
            margin-right: 8px;
            color: #555;
        }
        .footer {
            background-color: #f8f9fa;
            padding: 20px 30px;
            text-align: center;
            border-top: 1px solid #e1e5e9;
            color: #666;
            font-size: 13px;
        }
        .status-badge {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 12px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
        }
        .status-active {
            background-color: #d4edda;
            color: #155724;
        }
        .status-resolved {
            background-color: #d1ecf1;
            color: #0c5460;
        }
        .region-badge {
            display: inline-block;
            padding: 2px 6px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: 600;
            background-color: #e9ecef;
            color: #495057;
        }
    </style>
</head>
<body>
    <div class="container">
            <div class="header">
            <h1>预警通知汇总</h1>
            <p>生成时间: {{.GenerateTime}}</p>
        </div>
    
        <div class="content">
    <div class="summary">
                    <h3>统计摘要</h3>
        <ul>
                    <li>总预警数量: <strong>{{.TotalCount}}</strong></li>
            <li>涉及域名数量: <strong>{{.DomainCount}}</strong></li>
                    <li>预警来源: <strong>{{.SourceCount}}</strong> 个</li>
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
                        检测到域名<span class="domain-highlight">【{{$alert.Domain}}】</span>{{$alert.Message}}
                    </div>
                    <div class="alert-details">
                        <div class="detail-item">
                            <span class="detail-label">时间:</span>
                            <span>{{$alert.AlertTime.Format "2006-01-02 15:04:05"}}</span>
                        </div>
                        <div class="detail-item">
                            <span class="detail-label">来源:</span>
                            <span>{{$alert.Source}}</span>
                        </div>
                        <div class="detail-item">
                            <span class="detail-label">状态:</span>
                            <span class="status-badge {{if eq $alert.Status "active"}}status-active{{else}}status-resolved{{end}}">{{$alert.Status}}</span>
                        </div>
                        <div class="detail-item">
                            <span class="detail-label">区域:</span>
                            <span class="region-badge">{{$alert.Region}}</span>
                        </div>
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
		DomainCount  int
		SourceCount  int
		Alerts       []Alert
	}

	// 统计信息
	domainSet := make(map[string]bool)
	sourceSet := make(map[string]bool)
	
	for _, alert := range alerts {
		domainSet[alert.Domain] = true
		sourceSet[alert.Source] = true
	}

	// 确定时间范围
	var startTime, endTime time.Time
	if len(alerts) > 0 {
		startTime = alerts[len(alerts)-1].AlertTime
		endTime = alerts[0].AlertTime
	}

	data := TemplateData{
		GenerateTime: time.Now().Format("2006-01-02 15:04:05"),
		StartTime:    startTime.Format("2006-01-02 15:04:05"),
		EndTime:      endTime.Format("2006-01-02 15:04:05"),
		TotalCount:   len(alerts),
		DomainCount:  len(domainSet),
		SourceCount:  len(sourceSet),
		Alerts:       alerts,
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