package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

type Email struct {
	*SMTPInfo
}

// 定义 SMTPinfo 结构体用于传递发送邮件所必需的信息
type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

func NewEmail(info *SMTPInfo) *Email {
	return &Email{SMTPInfo: info}
}

func (e *Email) SendMail(to []string, subject, body string) error {
	// gomail.NewMessage() 创建一个消息实例
	m := gomail.NewMessage()
	// 设置邮件的一些必要信息
	m.SetHeader("From", e.From)     // 发件人
	m.SetHeader("To", to...)        // 收件人 to...: 将 to 切片打散
	m.SetHeader("Subject", subject) // 邮件主题
	m.SetBody("text/html", body)    // 邮件正文

	// gomail.NewDialer() 创建一个新的 SMTP 拨号实例，设置对应的拨号信息用于连接 SMTP 服务器
	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	// DialAndSend() 打开与 SMTP 服务器的连接并发送电子邮件
	return dialer.DialAndSend(m)
}
