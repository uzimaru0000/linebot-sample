package main

import (
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

func Callback(w http.ResponseWriter, req *http.Request) {
	events, err := client.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
	}

	for _, event := range events {
		log.Printf("Got event %v", event) //log of event
		switch event.Type {
		case linebot.EventTypeMessage:

		case linebot.EventTypePostback:
		}
	}
}
