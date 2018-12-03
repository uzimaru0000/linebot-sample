package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/poccariswet/bot/auth"
	"github.com/poccariswet/bot/template"
)

var client *linebot.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	client, err = auth.NewBot(os.Getenv("LINE_SECRET"), os.Getenv("LINE_TOKEN"))
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	http.HandleFunc("/callback", Callback)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func Callback(w http.ResponseWriter, req *http.Request) {
	events, err := client.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			log.Printf("bat request")
			w.WriteHeader(400)
		} else {
			log.Printf("server error")
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		log.Printf("Got event %v", event) //log of event
		switch event.Type {
		case linebot.EventTypeMessage:
			switch msg := event.Message.(type) {
			case *linebot.TextMessage:
				log.Printf("%v", msg)
				if err = MessageHander(msg, event.ReplyToken); err != nil {
					log.Print(err)
				}
			}
		case linebot.EventTypePostback:

		}
	}
}

func MessageHander(message *linebot.TextMessage, token string) error {
	var msg linebot.SendingMessage
	switch message.Text {
	case "buttons":
		btn := template.NewButtons()
		if err := btn.AddButtons(
			linebot.NewPostbackAction("Say hello1", "hello こんにちは", "", "hello こんにちは"),
			linebot.NewPostbackAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
			linebot.NewPostbackAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
			linebot.NewPostbackAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
		); err != nil {
			return err
		}
		msg = btn.ButtonsTemplate()

	case "confirm":
		confirm := template.NewConfirms()
		msg = confirm.ConfirmsTemplate()

	case "carousel":
		carousel := template.NewCarousel()
		btn := template.NewButtons()
		if err := btn.AddButtons(
			linebot.NewPostbackAction("Say hello1", "hello こんにちは", "", "hello こんにちは"),
			linebot.NewPostbackAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
			linebot.NewPostbackAction("言 hello2", "hello こんにちは", "hello こんにちは", ""),
		); err != nil {
			return err
		}

		if err := carousel.SetColumns(btn, btn); err != nil {
			return err
		}
		log.Println(carousel.Columns)
		msg = carousel.CarouselTemplate()

	case "image carousel":
		col := template.NewImageColumns()
		col.SetImageAction(linebot.NewURIAction("Go to LINE", "https://line.me"))

		c := template.NewImageCarousel()
		if err := c.SetImageCarousel(col, col, col); err != nil {
			return err
		}
		msg = c.CarouselTemplate()

	default:
		log.Println(message.Text)
	}
	if msg != nil {
		if _, err := client.ReplyMessage(token, msg).Do(); err != nil {
			return err
		}
	}
	return nil
}
