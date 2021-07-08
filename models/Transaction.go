package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kangata/gonotes/internal/date"
	"github.com/kangata/gonotes/internal/env"
	"github.com/kangata/gonotes/internal/telegram"
)

type Transaction struct {
	Date        string      `json:"date"`
	Customer    string      `json:"customer"`
	Product     string      `json:"product"`
	Price       uint64      `json:"price"`
	WebhookBody WebhookBody `json:"webhhok_body"`
}

func (t *Transaction) Parse(body *WebhookBody) Message {
	text := strings.Split(body.Message.Text, "#")

	if len(text) < 4 {
		return t
	}

	t.Date = time.Unix(int64(body.Message.Date), 0).Format("01/02/06")

	if len(text) > 4 {
		if val, err := date.Parse(text[4]); err == nil {
			t.Date = val.Format("01/02/06")
		}
	}

	t.Customer = text[1]
	t.Product = text[2]

	if price, err := strconv.ParseUint(text[3], 10, 64); err != nil {
		fmt.Println(err)
	} else {
		t.Price = price
	}

	t.WebhookBody = *body

	return t
}

func (t *Transaction) ToSheetRow() []interface{} {
	return []interface{}{
		t.Date,
		t.Customer,
		t.Product,
		t.Price,
		t.WebhookBody.Message.From.FirstName + "@" + env.Get("APP_NAME"),
	}
}

func (t *Transaction) SendAsTelegramMessage() (*http.Response, error) {
	ID := t.WebhookBody.Message.Chat.ID
	text := fmt.Sprintf("%s %s %s %d", t.Date, t.Customer, t.Product, t.Price)

	return telegram.SendTextMessage(ID, text)
}
