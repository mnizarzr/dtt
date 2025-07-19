package utils

import (
	"bytes"
	"html/template"
	"path/filepath"

	"github.com/mnizarzr/dot-test/config"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.Config
}

type EmailData struct {
	Name     string
	Email    string
	AppName  string
	Subject  string
	Content  string
	Role     string
	Password string
}

// NewEmailService creates a new email service instance
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendEmail sends an email using the configured SMTP settings
func (e *EmailService) SendEmail(to, subject, htmlBody, textBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.SmtpFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", textBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewDialer(e.config.SmtpHost, e.config.SmtpPort, e.config.SmtpUser, e.config.SmtpPassword)

	return d.DialAndSend(m)
}

// RenderTemplate renders an HTML email template with the provided data
func (e *EmailService) RenderTemplate(templateName string, data EmailData) (string, error) {
	templatePath := filepath.Join("template", "email", templateName+".html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SendWelcomeEmail sends a welcome email to a new user
func (e *EmailService) SendWelcomeEmail(userEmail, userName, userRole, userPlainPassword string) error {
	data := EmailData{
		Name:    userName,
		Email:   userEmail,
		Role:	userRole,
		Password: userPlainPassword,
		AppName: e.config.AppName,
		Subject: "Welcome to " + e.config.AppName,
	}

	htmlBody, err := e.RenderTemplate("welcome", data)
	if err != nil {
		return err
	}

	// Simple text fallback
	textBody := "Welcome to " + e.config.AppName + "!\n\nHi " + userName + ",\n\nThank you for registering with us. We're excited to have you on board!\n\nBest regards,\nThe " + e.config.AppName + " Team"

	return e.SendEmail(userEmail, data.Subject, htmlBody, textBody)
}
