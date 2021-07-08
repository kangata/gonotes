package models

import "net/http"

type Message interface {
	Parse(body *WebhookBody) Message

	ToSheetRow() []interface{}

	SendAsTelegramMessage() (*http.Response, error)
}
