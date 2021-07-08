package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kangata/gonotes/internal/env"
	"github.com/kangata/gonotes/internal/spreadsheet"
	"github.com/kangata/gonotes/models"
)

var sheet spreadsheet.Spreadsheet

func getRequestIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")

	if forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr
}

func handle(res http.ResponseWriter, req *http.Request) {
	request := fmt.Sprintf("%s %s %s %s", getRequestIP(req), req.Method, req.URL, req.Header.Get("User-Agent"))

	log.Println(request)

	defer req.Body.Close()

	raw, err := ioutil.ReadAll(req.Body)

	if err == nil {
		log.Println(string(raw))
	}

	body := &models.WebhookBody{}

	if err := json.NewDecoder(bytes.NewBuffer(raw)).Decode(body); err != nil {
		log.Println("Could not decode request body", err)

		return
	}

	if !body.IsValidUser() {
		log.Printf("ID not in whitelist: %d\n", body.Message.From.ID)

		return
	}

	if val, err := json.Marshal(body); err == nil {
		log.Println(string(val))
	}

	message, err := body.ToMessage()

	if err != nil {
		log.Println(err)

		return
	}

	if val, err := json.Marshal(message); val != nil || err == nil {
		log.Println(string(val))
	}

	if _, err := sheet.AddItem(message); err != nil {
		log.Println(err)
	} else {
		if _, err := message.SendAsTelegramMessage(); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(err)
	}

	if _, err := sheet.NewService(); err != nil {
		log.Fatalln(err)
	}

	log.Println("App listen at http://localhost:" + env.Get("APP_PORT"))

	http.ListenAndServe(":"+env.Get("APP_PORT"), http.HandlerFunc(handle))
}
