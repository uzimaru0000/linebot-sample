package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/uzimaru0000/linebot/model"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/poccariswet/bot/auth"
	"github.com/poccariswet/bot/template"
)

var client *linebot.Client
var userData map[string]State
var appID string

type State = uint8

const (
	LISTEN = iota
	ASK
)

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

	appID = os.Getenv("APP_ID")

	userData = make(map[string]State)

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
				if err = MessageHandler(event.Source.UserID, msg, event.ReplyToken); err != nil {
					log.Print(err)
				}
			}
		case linebot.EventTypePostback:
			data := event.Postback.Data
			if err = PostBackHandler(event.Source.UserID, data, event.ReplyToken); err != nil {
				log.Print(err)
			}
		}
	}
}

func MessageHandler(userID string, message *linebot.TextMessage, token string) error {
	var msg linebot.SendingMessage
	switch userData[userID] {
	case ASK:
		categories, err := model.GetCategory(appID)
		if err != nil {
			return err
		}
		id := model.MatchCategory(categories, message.Text)
		recipes, err := model.GetRecipe(appID, id)
		if err != nil {
			return err
		}

		carousel := template.NewCarousel()
		btns := make([]*template.Buttons, 0)
		for _, recipe := range recipes {
			btns = append(btns, recipe.RecipeTemplate())
		}
		if err := carousel.SetColumns(btns...); err != nil {
			return err
		}
		msg = carousel.CarouselTemplate()

		userData[userID] = LISTEN
	default:
		btn := template.NewButtons()
		btn.Title = "レシピをおすすめします"
		btn.SubTitle = "3つから選んでください！"
		btn.ImagePath = "https://3.bp.blogspot.com/-N2OBmlrmp6I/UnyHSqHeW3I/AAAAAAAAahc/1XbLO4ZbaQg/s800/cooking_chef.png"
		if err := btn.AddButtons(
			linebot.NewPostbackAction("今の気分からオススメ", "1", "", "今の気分からオススメ"),
			linebot.NewPostbackAction("季節のオススメ", "2", "", "季節のオススメ"),
			linebot.NewPostbackAction("おまかせ！", "3", "", "おまかせ！"),
		); err != nil {
			return err
		}
		msg = btn.ButtonsTemplate()
	}

	if msg != nil {
		if _, err := client.ReplyMessage(token, msg).Do(); err != nil {
			return err
		}
	}
	return nil
}

func PostBackHandler(userID, data, token string) error {
	var msg linebot.SendingMessage
	switch data {
	case "1":
		msg = linebot.NewTextMessage("今の気分は？")
		userData[userID] = ASK
	case "2":
		id := model.TimeToCategoryID(time.Now())
		recipes, err := model.GetRecipe(appID, id)
		if err != nil {
			return err
		}

		carousel := template.NewCarousel()
		btns := make([]*template.Buttons, 0)
		for _, recipe := range recipes {
			btns = append(btns, recipe.RecipeTemplate())
		}
		if err := carousel.SetColumns(btns...); err != nil {
			return err
		}
		msg = carousel.CarouselTemplate()
	case "3":
		// TODO: ランダムなカテゴリーIDでおすすめ
	}

	if msg != nil {
		if _, err := client.ReplyMessage(token, msg).Do(); err != nil {
			return err
		}
	}
	return nil
}
