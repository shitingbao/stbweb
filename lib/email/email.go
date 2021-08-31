package email

import (
	"net/smtp"
	"strings"
)

var (
	// https://mail.qq.com/cgi-bin/frame_html?sid=8j93veUohGokHS5O&r=8de5618177b9f4b48ea3fb4855412278
	mailSmtpPas  = "xqeyhlzuzmjbbfha" // 邮箱的授权码，去邮箱自己获取
	mailSmtpPort = ":587"
	fromEmail    = "871258317@qq.com"
	mailSmtpHost = "smtp.qq.com" // 此处填写SMTP服务器
)

// SendEmail发送验证码
func SendEmail(username, title string, toEmail []string, content string) error {
	auth := smtp.PlainAuth("", fromEmail, mailSmtpPas, mailSmtpHost)
	contentType := "Content-Type: text/plain; charset=UTF-8"

	msg := []byte("To: " + strings.Join(toEmail, ",") + "\r\nFrom: " + username +
		"<" + fromEmail + ">\r\nSubject: " + title + "\r\n" + contentType + "\r\n\r\n" + "your code is:" + content)
	err := smtp.SendMail(mailSmtpHost+mailSmtpPort, auth, fromEmail, toEmail, msg)
	if err != nil {
		return err
	}
	return nil
}
