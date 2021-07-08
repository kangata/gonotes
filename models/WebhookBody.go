package models

import (
	"errors"
	"strconv"
	"strings"

	"github.com/kangata/gonotes/internal/env"
)

type WebhookBody struct {
	Message struct {
		From struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
		} `json:"from"`
		Date uint64 `json:"date"`
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func (body *WebhookBody) ToMessage() (Message, error) {
	text := strings.Split(body.Message.Text, "#")
	messageType := text[0]

	if len(text) < 1 {
		return nil, errors.New("message invalid format")
	}

	if strings.ToUpper(messageType) == "TRX" {
		return new(Transaction).Parse(body), nil
	} else if strings.ToUpper(messageType) == "PAY" {
		return new(Payment).Parse(body), nil
	}

	return nil, errors.New("message invalid type")
}

func (body *WebhookBody) IsValidUser() bool {
	whitelistID := strings.Split(env.Get("TELEGRAM_WHITELIST_ID"), ",")

	for _, id := range whitelistID {
		if val, err := strconv.ParseInt(id, 10, 64); err == nil {
			if val == body.Message.From.ID {
				return true
			}
		}
	}

	return false
}
