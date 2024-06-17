package server

import (
	"log"
	"os"
	"time"

	"gopkg.in/telebot.v3"
)

func NewBot() *telebot.Bot {
	pref := telebot.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	return b
}
