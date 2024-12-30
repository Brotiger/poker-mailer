package controller

import (
	"context"
	"fmt"

	"github.com/Brotiger/poker-mailer/internal/service"
	"github.com/nats-io/nats.go/jetstream"
)

type MailController struct {
	mailService *service.MailService
}

func NewMailController() (*MailController, error) {
	mailService, err := service.NewMailService()
	if err != nil {
		return nil, fmt.Errorf("failed to create mail service, error: %w", err)
	}

	return &MailController{
		mailService: mailService,
	}, nil
}

func (mc *MailController) Send(ctx context.Context, msg jetstream.Msg) error {
	message, err := mc.mailService.GetMessage(msg.Headers().Get("type"), msg.Data())
	if err != nil {
		return fmt.Errorf("failed to get message, error: %w", err)
	}

	if err := mc.mailService.Send(message, msg.Headers().Get("email")); err != nil {
		return fmt.Errorf("failed to send, error: %w", err)
	}

	return nil
}
