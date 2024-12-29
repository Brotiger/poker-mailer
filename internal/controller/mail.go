package controller

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type MailController struct{}

func NewMailController() *MailController {
	return &MailController{}
}

func (m *MailController) Send(ctx context.Context, msg jetstream.Msg) error {
	return nil
}
