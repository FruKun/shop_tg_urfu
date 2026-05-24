package main

import (
	"log"
	"net/http"

	"github.com/FruKun/shop_tg_bot_urfu/botupdates"
	"github.com/FruKun/shop_tg_bot_urfu/config"
	"github.com/FruKun/shop_tg_bot_urfu/db"
	"github.com/FruKun/shop_tg_bot_urfu/webadmin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.New()

	storage, err := db.New(cfg.DBPath)
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	handler := botupdates.New(bot, storage, cfg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			go handler.HandleUpdate(update)
		}
	}()

	adminHandler := webadmin.New(storage, cfg)
	mux := http.NewServeMux()
	adminHandler.RegisterRoutes(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
	})

	go func() {
		if err := http.ListenAndServe(cfg.WebPort, mux); err != nil {
			log.Panic(err)
		}
	}()
	select {}
}
