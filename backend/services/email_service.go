package services

import (
	"fmt"
	"net/smtp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EmailService struct {
	Host        string
	Port        string
	User        string
	Password    string
	From        string
	FrontendUrl string
	AppName     string
}

func NewEmailService(host, port, user, password, from, frontendUrl, appName string) *EmailService {
	return &EmailService{
		Host:        host,
		Port:        port,
		User:        user,
		Password:    password,
		From:        from,
		FrontendUrl: frontendUrl,
		AppName:     appName,
	}
}

func (s *EmailService) SendRecoveryEmail(toEmail, token string) error {
	recoveryLink := fmt.Sprintf("%s/reset-password?token=%s", s.FrontendUrl, token)

	appNameClean := strings.ReplaceAll(s.AppName, "-", " ")
	displayName := cases.Title(language.BrazilianPortuguese).String(appNameClean)
	fromHeader := fmt.Sprintf("From: \"%s\" <%s>\n", displayName, s.From)

	subject := "Subject: Recovery password\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<div style="font-family: sans-serif; background-color: #f9f9f9; padding: 20px; color: #333;">
        <div style="max-width: 600px; margin: 0 auto; background-color: #fff; padding: 30px; border-radius: 8px; border: 1px solid #eee;">
            <h2 style="color: #111; margin-top: 0;">Recuperação de Senha</h2>
            <p style="font-size: 16px; line-height: 1.5; color: #555;">Você solicitou a alteração de sua senha. Use o link abaixo para redefini-la:</p>
            
            <p style="margin: 30px 0; text-align: center;">
                <a href="%s" style="background-color: #007bff; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: bold; display: inline-block;">
                    Definir Nova Senha
                </a>
            </p>
            
            <hr style="border: 0; border-top: 1px solid #eee; margin: 30px 0;">
            <p style="font-size: 12px; color: #888; line-height: 1.5;">
                Este link expira em 15 minutos. Se você não solicitou essa alteração, pode ignorar este e-mail com segurança.
            </p>
        </div>
    </div>
	`, recoveryLink)

	message := []byte(fromHeader + subject + mime + body)
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	var auth smtp.Auth
	if s.User != "" && s.Password != "" {
		auth = smtp.PlainAuth("", s.User, s.Password, s.Host)
	}

	err := smtp.SendMail(addr, auth, s.From, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
