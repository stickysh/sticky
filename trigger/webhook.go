package trigger

import (
	ctrlID "github.com/stickysh/sticky/pkg/id"

)

type WebhookID string

type Webhook struct {
	ID        WebhookID
	Action    string
	Enabled   bool
}

func NewWebhook(action string) *Webhook {
	return &Webhook{
		ID: WebhookID(ctrlID.GenULID()),
		Action: action,
		Enabled: true,
	}
}

type WebhookRepo interface {
	Store(t *Webhook) error

	Find(id WebhookID) (*Webhook, error)
}