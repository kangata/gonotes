package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kangata/gonotes/internal/env"
)

var baseURL string = "https://api.telegram.org/bot"

type TextMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func SendTextMessage(ID int64, text string) (*http.Response, error) {
	message := TextMessage{
		ChatID: ID,
		Text:   text,
	}

	url := fmt.Sprintf("%s%s/%s", baseURL, env.Get("TELEGRAM_BOT_TOKEN"), "sendMessage")

	json, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	return http.Post(url, "application/json", bytes.NewBuffer(json))
}
