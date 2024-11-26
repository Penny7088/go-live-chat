package emailtool

import (
	"fmt"
	"net/smtp"

	"lingua_exchange/internal/config"
	"lingua_exchange/pkg/strutil"
)

func SendEmail(email string, code string, subject string, fileName string) error {
	userName := config.Get().SMTP.UserName
	password := config.Get().SMTP.Password
	host := config.Get().SMTP.Host
	port := config.Get().SMTP.Port
	// 设置 SMTP 服务器地址和端口
	addr := fmt.Sprintf("%s:%d", host, port)
	// 创建身份验证信息
	auth := smtp.PlainAuth("", userName, password, host)
	templatePath, err := strutil.GetTemplatePath(fileName)
	if err != nil {
		return err
	}
	renderTemplate, err := strutil.RenderTemplate(templatePath, code)
	if err != nil {
		return err
	}
	message := []byte(fmt.Sprintf("MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"Subject: %s\r\n\r\n%s", subject, renderTemplate))

	// 发送邮件
	err = smtp.SendMail(addr, auth, userName, []string{email}, message)
	return err
}
