package botupdates

import (
	"fmt"
	"log"
	"strings"

	"github.com/FruKun/shop_tg_bot_urfu/keyboard"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) help(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "help")
	msg.ReplyMarkup = keyboard.MainMenu()
	h.bot.Send(msg)
}

func (h *Handler) cart(message *tgbotapi.Message) {
	items, err := h.storage.GetCart(message.From.ID)
	var msg tgbotapi.MessageConfig
	msg.ReplyMarkup = keyboard.MainMenu()
	if err != nil {
		log.Println(err)
		msg = tgbotapi.NewMessage(message.Chat.ID, "error")
		h.bot.Send(msg)
		return
	}
	if len(items) == 0 {
		msg = tgbotapi.NewMessage(message.Chat.ID, "корзина пуста")
		h.bot.Send(msg)
		return
	}
	var total float64
	var text strings.Builder
	text.WriteString("Ваша корзина \n")

	for i, item := range items {
		product, err := h.storage.GetProduct(item.ProductID)
		if err != nil {
			log.Println(err)
			continue
		}
		subtotal := product.Price * float64(item.Quantity)
		total += subtotal
		fmt.Fprintf(&text, "%d. %s x%d = %.2f руб \n", i+1, product.Name, item.Quantity, subtotal)
	}
	fmt.Fprintf(&text, "\n Итого: %.2f руб ", total)
	msg = tgbotapi.NewMessage(message.Chat.ID, text.String())
	msg.ReplyMarkup = keyboard.CartMenu()
	h.bot.Send(msg)
}

func (h *Handler) catalog(message *tgbotapi.Message, page int) {
	products, err := h.storage.GetAllProduct()

	var msg tgbotapi.MessageConfig
	msg.ReplyMarkup = keyboard.MainMenu()

	if err != nil {
		log.Println(err)
		msg = tgbotapi.NewMessage(message.Chat.ID, "ошибка загрузки")
		h.bot.Send(msg)
		return
	}

	if len(products) == 0 {
		msg = tgbotapi.NewMessage(message.Chat.ID, "товаров нет")
		h.bot.Send(msg)
		return
	}

	itemsPerPage := 5
	totalPages := (len(products) + itemsPerPage - 1) / itemsPerPage

	start := page * itemsPerPage
	end := min(start+itemsPerPage, len(products))
	if start == end {
		start = max(start-itemsPerPage, 0)
		page = max(page-1, 0)
	}

	var text strings.Builder
	for i, product := range products[start:end] {
		fmt.Fprintf(&text, "%d. %s\nЦена: %.2f руб.\n%s\nВ наличии: %d\n\n",
			i+1, product.Name, product.Price, product.Description, product.Quantity)
	}
	fmt.Fprintf(&text, "сраница %d из %d\n", page+1, totalPages)

	msg = tgbotapi.NewMessage(message.Chat.ID, text.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard.CatalogMenu(page, totalPages, products)
	h.bot.Send(msg)
}
