package server

import (
	"log"
	"os"

	"gopkg.in/telebot.v3"
)

func NewBot() *telebot.Bot {
	pref := telebot.Settings{
		Token: os.Getenv("BOT_TOKEN"),
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	return b
}
