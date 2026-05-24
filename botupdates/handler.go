package botupdates

import (
	"strings"

	"github.com/FruKun/shop_tg_bot_urfu/config"
	"github.com/FruKun/shop_tg_bot_urfu/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot     *tgbotapi.BotAPI
	storage *db.Database
	config  *config.Config
}

func New(bot *tgbotapi.BotAPI, storage *db.Database, cfg *config.Config) *Handler {
	return &Handler{
		bot:     bot,
		storage: storage,
		config:  cfg,
	}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		h.handleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		h.handleCallback(update.CallbackQuery)
	}
}
func (h *Handler) handleMessage(message *tgbotapi.Message) {
	switch strings.ToLower(message.Text) {
	case "/catalog", "каталог":
		h.catalog(message, 0)
	case "/cart", "корзина":
		h.cart(message)
	case "/help", "/start":
		h.help(message)
	default:
		h.help(message)
	}
}

func (h *Handler) handleCallback(callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	switch {
	case strings.HasPrefix(data, "catalog_page_"):
		h.catalogPage(callback)
	case strings.HasPrefix(data, "cart_add_"):
		h.cart_add(callback)
	case strings.HasPrefix(data, "product_description_"):
		h.product_description(callback)
	case strings.HasPrefix(data, "cart_menu_order"):
		h.cart_menu_order(callback)
	case strings.HasPrefix(data, "cart_menu_clear_all"):
		h.cart_menu_clear_all(callback)
	}

}
