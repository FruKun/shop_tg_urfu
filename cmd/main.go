package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/FruKun/shop_tg_bot_urfu/botupdates"
	"github.com/FruKun/shop_tg_bot_urfu/config"
	"github.com/FruKun/shop_tg_bot_urfu/db"
	"github.com/FruKun/shop_tg_bot_urfu/logger"
	"github.com/FruKun/shop_tg_bot_urfu/webadmin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.New()

	writer := io.Writer(os.Stdout)
	if cfg.LogDir != "" {
		file, err := os.OpenFile(cfg.LogDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Panic(err)
		}
		writer = io.MultiWriter(writer, file)
		defer file.Close()
	}

	loglvl := "DEBUG"
	if cfg.LogLevel != "" {
		loglvl = cfg.LogLevel
	}

	baseLogger := logger.New(writer, loglvl)
	mainLogger := baseLogger.NewModule("main")

	mainLogger.Info("start database dbpath: %s", cfg.DBPath)
	storage, err := db.New(cfg.DBPath, baseLogger.NewModule("db"))
	if err != nil {
		mainLogger.Panic("%s", err)
	}

	mainLogger.Info("start bot, Bot Token: %s", string([]rune(cfg.BotToken)[len(cfg.BotToken)-5:]))
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		mainLogger.Panic("%s", err)
	}

	bot.Debug = false
	if strings.ToUpper(loglvl) == "DEBUG" {
		bot.Debug = true
	}

	mainLogger.Info("start handle update")
	handler := botupdates.New(bot, storage, cfg, baseLogger.NewModule("botupdates"))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			go handler.HandleUpdate(update)
		}
	}()

	adminHandler := webadmin.New(storage, cfg, baseLogger.NewModule("webadmin"))
	mux := http.NewServeMux()
	adminHandler.RegisterRoutes(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
	})

	mainLogger.Info("start web")
	go func() {
		if err := http.ListenAndServe(cfg.WebPort, mux); err != nil {
			mainLogger.Panic("%s", err)
		}
	}()
	select {}
}
