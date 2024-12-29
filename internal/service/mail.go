package service

import (
	"bytes"
	"fmt"
	"html/template"

	natsModel "github.com/Brotiger/poker-core_api/pkg/nats/model"
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

func (ms *MailService) GetMessage(messageType string, data []byte) error {
	var body bytes.Buffer

	template := ms.templateMap[messageType]
	switch messageType {
	case "register":
		modelRegister := natsModel.Register{}
		if err := template.Execute(&body, modelRegister); err != nil {
			return fmt.Errorf("failed to tempalte execute, error: %w", err)
		}
	case "restore":
		modelRestore := natsModel.Restore{}
		if err := template.Execute(&body, modelRestore); err != nil {
			return fmt.Errorf("failed to tempalte execute, error: %w", err)
		}
	}

	return nil
}
