package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// EmailConfig 邮件配置结构
type EmailConfig struct {
	APIUrl      string   `json:"api_url"`
	AppID       string   `json:"app_id"`
	AppSecret   string   `json:"app_secret"`
	From        string   `json:"from"`
	To          []string `json:"to"`
	DebugMode   bool     `json:"debug_mode"`
	DebugAPIUrl string   `json:"debug_api_url"`
}

// UserInfo 用户信息结构
type UserInfo struct {
	Name  string `json:"name"`
	EName string `json:"e_name"`
	Email string `json:"email"`
}

// EmailAPIRequest 邮件API请求结构
type EmailAPIRequest struct {
	OpdAppid     string `json:"opdAppid"`
	OpdAppsecret string `json:"opdAppsecret"`
	ToList       string `json:"to_list"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	Mimetype     string `json:"mimetype"`
}

// EmailAPIResponse 邮件API响应结构
type EmailAPIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// RecipientInfo 收件人信息
type RecipientInfo struct {
	Email string
	Found bool
}

var emailConfig EmailConfig
var userList []UserInfo

// InitEmailConfig 初始化邮件配置
func InitEmailConfig() {
	emailConfig = config.Email

	if err := loadUserList(); err != nil {
		LogSystem(logrus.ErrorLevel, "email", "加载用户列表失败", map[string]interface{}{
			"error": err.Error(),
		})
		log.Printf("加载用户列表失败: %v", err)
	}
}

// loadUserList 加载用户列表
func loadUserList() error {
	file, err := os.Open("userlist.json")
	if err != nil {
		return fmt.Errorf("打开用户列表文件失败: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&userList); err != nil {
		return fmt.Errorf("解析用户列表JSON失败: %v", err)
	}

	LogSystem(logrus.InfoLevel, "email", "用户列表加载成功", map[string]interface{}{
		"user_count": len(userList),
	})
	log.Printf("用户列表加载成功，共 %d 个用户", len(userList))

	return nil
}

// findUserByEName 根据英文名查找用户邮箱
func findUserByEName(eName string) (string, bool) {
	for _, user := range userList {
		if user.EName == eName && user.Email != "" {
			return user.Email, true
		}
	}
	return "", false
}

// SendAlertEmail 发送预警通知邮件（按用户分组）
func SendAlertEmail(userAlertsList []UserAlerts) error {
	if len(userAlertsList) == 0 {
		LogSystem(logrus.WarnLevel, "email", "没有预警信息需要发送", nil)
		return fmt.Errorf("没有预警信息需要发送")
	}

	LogSystem(logrus.InfoLevel, "email", "开始发送邮件", map[string]interface{}{
		"user_count": len(userAlertsList),
	})

	var successCount, failCount int
	var successRecipients, failRecipients []string
	var notFoundUsers []string
	var fallbackAlerts []UserAlerts
	var fallbackEmailSent bool

	for _, userAlerts := range userAlertsList {
		recipientInfo := generateRecipientEmail(userAlerts.Recipient)

		if !recipientInfo.Found {
			notFoundUsers = append(notFoundUsers, userAlerts.Recipient)
			fallbackAlerts = append(fallbackAlerts, userAlerts)
			continue
		}

		LogSystem(logrus.InfoLevel, "email", "准备发送用户邮件", map[string]interface{}{
			"recipient": userAlerts.Recipient,
			"email": recipientInfo.Email,
			"found": recipientInfo.Found,
			"alert_count": len(userAlerts.Alerts),
		})

		if err := sendEmailToUser(userAlerts, recipientInfo); err != nil {
			LogEmail(recipientInfo.Email, "预警通知", false, err.Error())
			log.Printf("发送邮件给用户 %s (%s) 失败: %v", 
				userAlerts.Recipient, recipientInfo.Email, err)
			failCount++
			failRecipients = append(failRecipients, recipientInfo.Email)
			continue
		}

		LogEmail(recipientInfo.Email, "预警通知", true, "")
		log.Printf("成功发送邮件给用户: %s (%s)，包含 %d 条预警信息", 
			userAlerts.Recipient, recipientInfo.Email, len(userAlerts.Alerts))
		successCount++
		successRecipients = append(successRecipients, recipientInfo.Email)
	}

	if len(fallbackAlerts) > 0 {
		fallbackEmail := "liyongchang@kugou.net"
		LogSystem(logrus.InfoLevel, "email", "发送合并管理员邮件", map[string]interface{}{
			"fallback_email": fallbackEmail,
			"not_found_users": notFoundUsers,
			"alert_count": len(fallbackAlerts),
		})

		if err := sendFallbackEmail(fallbackAlerts, notFoundUsers); err != nil {
			LogEmail(fallbackEmail, "管理员预警通知", false, err.Error())
			log.Printf("发送管理员邮件失败: %v", err)
			failCount++
			failRecipients = append(failRecipients, fallbackEmail)
		} else {
			LogEmail(fallbackEmail, "管理员预警通知", true, "")
			log.Printf("成功发送管理员邮件给: %s，包含 %d 个未找到用户的预警信息", 
				fallbackEmail, len(fallbackAlerts))
			successCount++
			successRecipients = append(successRecipients, fallbackEmail+"(管理员)")
			fallbackEmailSent = true
		}
	}

	LogSystem(logrus.InfoLevel, "email", "邮件发送完成", map[string]interface{}{
		"total_users": len(userAlertsList),
		"success_count": successCount,
		"fail_count": failCount,
		"not_found_count": len(notFoundUsers),
		"success_recipients": successRecipients,
		"fail_recipients": failRecipients,
		"not_found_users": notFoundUsers,
		"fallback_email_sent": fallbackEmailSent,
	})

	log.Printf("邮件发送总结")
	log.Printf("  成功: %d 个用户", successCount)
	log.Printf("  失败: %d 个用户", failCount)
	log.Printf("  未找到: %d 个用户", len(notFoundUsers))
	if len(successRecipients) > 0 {
		log.Printf("  成功收件人: %v", successRecipients)
	}
	if len(failRecipients) > 0 {
		log.Printf("  失败收件人: %v", failRecipients)
	}
	if len(notFoundUsers) > 0 {
		log.Printf("  未找到用户: %v", notFoundUsers)
	}

	if failCount > 0 {
		return fmt.Errorf("部分邮件发送失败，成功: %d，失败: %d", successCount, failCount)
	}

	return nil
}

// sendEmailToUser 发送邮件给特定用户
func sendEmailToUser(userAlerts UserAlerts, recipientInfo RecipientInfo) error {
	subject, body, err := generateEmailContentForUser(userAlerts, recipientInfo)
	if err != nil {
		return fmt.Errorf("生成邮件内容失败: %v", err)
	}

	if err := sendEmailViaAPI([]string{recipientInfo.Email}, subject, body); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

// generateRecipientEmail 生成收件人信息，包括邮箱地址
func generateRecipientEmail(recipient string) RecipientInfo {
	if strings.Contains(recipient, "@") {
		return RecipientInfo{
			Email: recipient,
			Found: true,
		}
	}

	if email, found := findUserByEName(recipient); found {
		LogSystem(logrus.InfoLevel, "email", "找到用户邮箱", map[string]interface{}{
			"e_name": recipient,
			"email": email,
		})
		return RecipientInfo{
			Email: email,
			Found: true,
		}
	}

	fallbackEmail := "liyongchang@kugou.net"
	LogSystem(logrus.WarnLevel, "email", "未找到用户，使用管理员邮箱", map[string]interface{}{
		"e_name": recipient,
		"fallback_email": fallbackEmail,
	})
	log.Printf("未找到用户 %s 的邮箱，使用管理员邮箱: %s", recipient, fallbackEmail)

	return RecipientInfo{
		Email: fallbackEmail,
		Found: false,
	}
}

// generateEmailContentForUser 为用户生成邮件内容
func generateEmailContentForUser(userAlerts UserAlerts, recipientInfo RecipientInfo) (string, string, error) {
	subject := fmt.Sprintf("预警通知 - %s - %s", userAlerts.Recipient, time.Now().Format("2006-01-02"))

	if !recipientInfo.Found {
		subject = fmt.Sprintf("【管理员】预警通知 - %s (未找到用户 - %s)", userAlerts.Recipient, time.Now().Format("2006-01-02"))
	}

	const emailTemplate = `<!DOCTYPE html>
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
            color: #007bff;
            font-size: 18px;
        }
        .summary p {
            margin: 5px 0;
            color: #6c757d;
        }
        .user-not-found {
            background-color: #f8d7da;
            border: 1px solid #f5c6cb;
            border-radius: 8px;
            padding: 15px;
            margin-bottom: 20px;
            color: #721c24;
        }
        .user-not-found h4 {
            margin: 0 0 10px 0;
            font-size: 16px;
        }
        .user-not-found ul {
            margin: 10px 0 0 0;
            padding-left: 20px;
        }
        .user-not-found li {
            margin-bottom: 5px;
        }
        .alert-item {
            background-color: #ffffff;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 15px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .alert-item h4 {
            margin: 0 0 10px 0;
            color: #dc3545;
            font-size: 16px;
        }
        .alert-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 15px;
            padding-top: 15px;
            border-top: 1px solid #e9ecef;
            font-size: 12px;
            color: #6c757d;
        }
        .alert-time {
            font-weight: 600;
        }
        .footer {
            background-color: #f8f9fa;
            padding: 20px 30px;
            text-align: center;
            color: #6c757d;
            font-size: 12px;
            border-top: 1px solid #e9ecef;
        }
        .footer p {
            margin: 5px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>预警通知</h1>
            <p>生成时间: {{.GenerateTime}}</p>
        </div>
    
        <div class="content">
            {{if not .UserFound}}
            <div class="user-not-found">
                <h4>用户信息提示</h4>
                <p>用户 <strong>{{.Recipient}}</strong> 在系统中未找到对应的邮箱地址，邮件已发送给系统管理员。</p>
                <p>请检查以下事项：</p>
                <ul>
                    <li>确认用户英文名拼写是否正确</li>
                    <li>检查用户列表是否包含该用户</li>
                    <li>确认用户邮箱地址是否有效</li>
                    <li>如需添加，请更新用户列表文件</li>
                </ul>
            </div>
            {{end}}
            
            <div class="summary">
                <h3>预警信息概览</h3>
                <p><strong>收件人:</strong> {{.Recipient}}</p>
                <p><strong>时间范围:</strong> {{.StartTime}} 至 {{.EndTime}}</p>
                <p><strong>预警数量:</strong> {{.TotalCount}} 条</p>
            </div>

            <h3 style="color: #dc3545; margin-bottom: 20px; font-size: 18px;">详细预警信息</h3>
            
            {{range $index, $alert := .Alerts}}
            <div class="alert-item">
                <h4>预警 #{{add $index 1}}</h4>
                <p style="margin: 10px 0; line-height: 1.6;">{{$alert.Message}}</p>
                <div class="alert-meta">
                    <span class="alert-time">时间: {{$alert.AlertTime.Format "2006-01-02 15:04:05"}}</span>
                </div>
            </div>
            {{end}}
        </div>
    
        <div class="footer">
            <p><strong>系统自动发送</strong> | 请及时处理相关预警信息</p>
            <p>如有疑问，请联系系统管理员</p>
        </div>
    </div>
</body>
</html>`

	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	if len(userAlerts.Alerts) > 0 {
		startTime = userAlerts.Alerts[0].AlertTime
		endTime = userAlerts.Alerts[0].AlertTime
		for _, alert := range userAlerts.Alerts {
			if alert.AlertTime.Before(startTime) {
				startTime = alert.AlertTime
			}
			if alert.AlertTime.After(endTime) {
				endTime = alert.AlertTime
			}
		}
	}

	type TemplateData struct {
		GenerateTime string
		StartTime    string
		EndTime      string
		TotalCount   int
		Recipient    string
		UserFound    bool
		Alerts       []Alert
	}

	data := TemplateData{
		GenerateTime: time.Now().Format("2006-01-02 15:04:05"),
		StartTime:    startTime.Format("2006-01-02 15:04:05"),
		EndTime:      endTime.Format("2006-01-02 15:04:05"),
		TotalCount:   len(userAlerts.Alerts),
		Recipient:    userAlerts.Recipient,
		UserFound:    recipientInfo.Found,
		Alerts:       userAlerts.Alerts,
	}

	tmpl, err := template.New("email").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}).Parse(emailTemplate)
	if err != nil {
		return "", "", fmt.Errorf("解析邮件模板失败: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("执行邮件模板失败: %v", err)
	}

	return subject, buf.String(), nil
}

// sendEmailViaAPI 通过HTTP API发送邮件
func sendEmailViaAPI(toUsers []string, subject, content string) error {
	apiURL := emailConfig.APIUrl
	if emailConfig.DebugMode && emailConfig.DebugAPIUrl != "" {
		apiURL = emailConfig.DebugAPIUrl
		log.Printf("使用调试模式邮件API: %s", apiURL)
	}

	toList := strings.Join(toUsers, ",")

	formData := url.Values{}
	formData.Set("opdAppid", emailConfig.AppID)
	formData.Set("opdAppsecret", emailConfig.AppSecret)
	formData.Set("to_list", toList)
	formData.Set("subject", subject)
	formData.Set("body", content)
	formData.Set("mimetype", "html")

	postData := formData.Encode()

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(postData))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Alert-System/1.0")
	req.Header.Set("Accept-Charset", "UTF-8")

	log.Printf("发送邮件信息")
	log.Printf("   收件人: %s", toList)
	log.Printf("   主题: %s", subject)
	log.Printf("   API地址: %s", apiURL)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	log.Printf("邮件API响应:")
	log.Printf("  状态码: %s", resp.Status)
	log.Printf("  响应内容: %s", string(respBody))

	var apiResp EmailAPIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		log.Printf("邮件API响应解析失败，原始响应: %s", string(respBody))
		return fmt.Errorf("解析API响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("邮件API返回错误状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	if apiResp.Code != 0 && apiResp.Code != 200 {
		return fmt.Errorf("邮件发送失败, code=%d, message=%s", apiResp.Code, apiResp.Message)
	}

	log.Printf("邮件发送成功: %s", apiResp.Message)
	return nil
}

// sendFallbackEmail 发送合并的管理员邮件
func sendFallbackEmail(fallbackAlerts []UserAlerts, notFoundUsers []string) error {
	fallbackRecipientInfo := RecipientInfo{
		Email: "liyongchang@kugou.net",
		Found: false,
	}

	subject, body, err := generateFallbackEmailContent(fallbackAlerts, notFoundUsers)
	if err != nil {
		return fmt.Errorf("生成管理员邮件内容失败: %v", err)
	}

	if err := sendEmailViaAPI([]string{fallbackRecipientInfo.Email}, subject, body); err != nil {
		return fmt.Errorf("发送管理员邮件失败: %v", err)
	}

	return nil
}

// generateFallbackEmailContent 生成管理员邮件内容（包含用户分组）
func generateFallbackEmailContent(fallbackAlerts []UserAlerts, notFoundUsers []string) (string, string, error) {
	subject := fmt.Sprintf("【管理员】预警通知 - %s (未找到用户) - %s", 
		strings.Join(notFoundUsers, ", "), time.Now().Format("2006-01-02"))

	const fallbackEmailTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>管理员预警通知</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; 
            margin: 0; 
            padding: 20px; 
            background-color: #f5f5f5; 
            line-height: 1.6;
        }
        .container {
            max-width: 900px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #dc3545 0%, #c82333 100%);
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
        .warning-box {
            background-color: #f8d7da;
            border: 1px solid #f5c6cb;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 30px;
            color: #721c24;
        }
        .warning-box h3 {
            margin: 0 0 15px 0;
            font-size: 18px;
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
            color: #007bff;
            font-size: 18px;
        }
        .summary p {
            margin: 5px 0;
            color: #6c757d;
        }
        .user-section {
            background-color: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .user-header {
            background-color: #007bff;
            color: white;
            padding: 15px;
            margin: -20px -20px 20px -20px;
            border-radius: 8px 8px 0 0;
        }
        .user-header h3 {
            margin: 0;
            font-size: 16px;
        }
        .alert-item {
            background-color: #ffffff;
            border: 1px solid #e9ecef;
            border-radius: 6px;
            padding: 15px;
            margin-bottom: 10px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .alert-item h4 {
            margin: 0 0 10px 0;
            color: #dc3545;
            font-size: 14px;
        }
        .alert-message {
            margin: 10px 0;
            line-height: 1.6;
        }
        .alert-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 10px;
            padding-top: 10px;
            border-top: 1px solid #e9ecef;
            font-size: 12px;
            color: #6c757d;
        }
        .alert-time {
            font-weight: 600;
        }
        .footer {
            background-color: #f8f9fa;
            padding: 20px 30px;
            text-align: center;
            color: #6c757d;
            font-size: 12px;
            border-top: 1px solid #e9ecef;
        }
        .footer p {
            margin: 5px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>【管理员】预警通知</h1>
            <p>未找到用户转发通知 | 生成时间: {{.GenerateTime}}</p>
        </div>
        
        <div class="content">
            <div class="warning-box">
                <h3>重要通知</h3>
                <p>系统在处理预警信息时，发现以下用户英文名在用户列表中未找到对应的邮箱地址。</p>
                <p><strong>未找到用户：</strong>{{.NotFoundUsers}}</p>
                <p>该邮件已发送给系统管理员，请及时处理相关用户信息。</p>
            </div>
            
            <div class="summary">
                <h3>预警信息概览</h3>
                <p><strong>未找到用户数量：</strong>{{.UserCount}} 个</p>
                <p><strong>总预警数量：</strong>{{.TotalAlerts}} 条</p>
                <p><strong>时间范围：</strong>{{.StartTime}} 至 {{.EndTime}}</p>
            </div>
            
            <h3 style="color: #dc3545; margin-bottom: 20px; font-size: 18px;">各用户详细预警信息</h3>
            
            {{range $userIndex, $userAlerts := .UserAlertsList}}
            <div class="user-section">
                <div class="user-header">
                    <h3>用户：{{$userAlerts.Recipient}} (未找到邮箱)</h3>
                </div>
                
                {{range $alertIndex, $alert := $userAlerts.Alerts}}
                <div class="alert-item">
                    <h4>预警 #{{add $alertIndex 1}}</h4>
                    <div class="alert-message">{{$alert.Message}}</div>
                    <div class="alert-meta">
                        <span class="alert-time">{{$alert.AlertTime.Format "2006-01-02 15:04:05"}}</span>
                    </div>
                </div>
                {{end}}
            </div>
            {{end}}
        </div>
        
        <div class="footer">
            <p><strong>系统自动发送给管理员</strong> | 请及时处理相关用户信息和预警</p>
            <p>建议：更新用户列表文件，添加缺失用户的邮箱地址</p>
        </div>
    </div>
</body>
</html>`

	var startTime, endTime time.Time
	totalAlerts := 0

	if len(fallbackAlerts) > 0 {
		startTime = fallbackAlerts[0].Alerts[0].AlertTime
		endTime = fallbackAlerts[0].Alerts[0].AlertTime

		for _, userAlerts := range fallbackAlerts {
			totalAlerts += len(userAlerts.Alerts)
			for _, alert := range userAlerts.Alerts {
				if alert.AlertTime.Before(startTime) {
					startTime = alert.AlertTime
				}
				if alert.AlertTime.After(endTime) {
					endTime = alert.AlertTime
				}
			}
		}
	}

	type FallbackTemplateData struct {
		GenerateTime   string
		NotFoundUsers  string
		UserCount      int
		TotalAlerts    int
		StartTime      string
		EndTime        string
		UserAlertsList []UserAlerts
	}

	data := FallbackTemplateData{
		GenerateTime:   time.Now().Format("2006-01-02 15:04:05"),
		NotFoundUsers:  strings.Join(notFoundUsers, ", "),
		UserCount:      len(notFoundUsers),
		TotalAlerts:    totalAlerts,
		StartTime:      startTime.Format("2006-01-02 15:04:05"),
		EndTime:        endTime.Format("2006-01-02 15:04:05"),
		UserAlertsList: fallbackAlerts,
	}

	tmpl, err := template.New("fallback_email").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}).Parse(fallbackEmailTemplate)
	if err != nil {
		return "", "", fmt.Errorf("解析管理员邮件模板失败: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("执行管理员邮件模板失败: %v", err)
	}

	return subject, buf.String(), nil
} 