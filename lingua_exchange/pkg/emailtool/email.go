package emailtool

import (
	"gopkg.in/gomail.v2"
	"lingua_exchange/internal/config"
	"lingua_exchange/pkg/strutil"
)

func CreatDialer() *gomail.Dialer {
	userName := config.Get().SMTP.UserName
	password := config.Get().SMTP.Password
	host := config.Get().SMTP.Host
	port := config.Get().SMTP.Port
	dialer := gomail.NewDialer(host, port, userName, password)
	dialer.SSL = true
	return dialer
}

func SendEmail(email string, code string, subject string, fileName string) error {
	dialer := CreatDialer()
	m := gomail.NewMessage()
	m.SetHeader("From", dialer.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)

	templatePath, err := strutil.GetTemplatePath(fileName)
	if err != nil {
		return err
	}
	renderTemplate, err := strutil.RenderTemplate(templatePath, code)
	if err != nil {
		return err
	}
	m.SetBody("text/html", renderTemplate)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
