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

type Payment struct {
	Date        string      `json:"date"`
	Customer    string      `json:"customer"`
	Method      string      `json:"method"`
	Amount      int64       `json:"amount"`
	WebhookBody WebhookBody `json:"webhhok_body"`
}

func (p *Payment) Parse(body *WebhookBody) Message {
	text := strings.Split(body.Message.Text, "#")

	if len(text) < 4 {
		return p
	}

	p.Date = time.Unix(int64(body.Message.Date), 0).Format("01/02/06")

	if len(text) > 4 {
		if val, err := date.Parse(text[4]); err == nil {
			p.Date = val.Format("01/02/06")
		}
	}

	p.Customer = text[1]
	p.Method = text[2]

	if amount, err := strconv.ParseInt(text[3], 10, 64); err != nil {
		fmt.Println(err)
	} else {
		p.Amount = amount * -1
	}

	p.WebhookBody = *body

	return p
}

func (p *Payment) ToSheetRow() []interface{} {
	return []interface{}{
		p.Date,
		p.Customer,
		fmt.Sprintf("PAY-%s", p.Method),
		p.Amount,
		p.WebhookBody.Message.From.FirstName + "@" + env.Get("APP_NAME"),
	}
}

func (p *Payment) SendAsTelegramMessage() (*http.Response, error) {
	ID := p.WebhookBody.Message.Chat.ID
	text := fmt.Sprintf("%s %s PAY-%s %d", p.Date, p.Customer, p.Method, p.Amount)

	return telegram.SendTextMessage(ID, text)
}
