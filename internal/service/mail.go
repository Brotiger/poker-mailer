package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"

	natsModel "github.com/Brotiger/poker-core_api/pkg/nats/model"
	"github.com/Brotiger/poker-mailer/internal/config"
)

type MailService struct {
	templateMap map[string]*template.Template
}

func NewMailService() (*MailService, error) {
	registerTemplate, err := template.ParseFiles("template/register.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse register template, error: %w", err)
	}

	restoreTemplate, err := template.ParseFiles("template/restore.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse restore template, error: %w", err)
	}

	return &MailService{
		templateMap: map[string]*template.Template{
			"register": registerTemplate,
			"restore":  restoreTemplate,
		},
	}, nil
}

func (ms *MailService) GetMessage(messageType string, data []byte) (string, error) {
	var body bytes.Buffer

	template := ms.templateMap[messageType]
	switch messageType {
	case "register":
		var modelRegister natsModel.Register
		if err := json.Unmarshal(data, &modelRegister); err != nil {
			return "", fmt.Errorf("failed to unmarshal, error: %w", err)
		}

		if err := template.Execute(&body, modelRegister); err != nil {
			return "", fmt.Errorf("failed to tempalte execute, error: %w", err)
		}
	case "restore":
		var modelRestore natsModel.Restore
		if err := json.Unmarshal(data, &modelRestore); err != nil {
			return "", fmt.Errorf("failed to unmarshal, error: %w", err)
		}

		if err := template.Execute(&body, modelRestore); err != nil {
			return "", fmt.Errorf("failed to tempalte execute, error: %w", err)
		}
	}

	message := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Subject: Приветственное письмо\r\n\r\n" +
		body.String()

	return message, nil
}

func (ms *MailService) Send(message, to string) error {
	auth := smtp.PlainAuth("", config.Cfg.SMTP.User, config.Cfg.SMTP.Password, config.Cfg.SMTP.Host)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", config.Cfg.SMTP.Host, config.Cfg.SMTP.Port), auth, config.Cfg.SMTP.From, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send mail, error: %w", err)
	}

	return nil
}
