package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/poccariswet/bot/auth"
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
