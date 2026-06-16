package botupdates

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/FruKun/shop_tg_bot_urfu/keyboard"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) callbackAnswer(callback *tgbotapi.CallbackQuery, str string) {
	if _, err := h.bot.Request(tgbotapi.NewCallback(callback.ID, str)); err != nil {
		h.logger.Error("%s", err)
	}

}

func (h *Handler) catalogPage(callback *tgbotapi.CallbackQuery) {
	var msg tgbotapi.MessageConfig
	msg.ReplyMarkup = keyboard.MainMenu()

	page, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "catalog_page_"))

	itemsPerPage := 5
	products, total, err := h.storage.GetPaginatedProducts(itemsPerPage, page)
	if err != nil {
		msg = tgbotapi.NewMessage(callback.Message.Chat.ID, err.Error())
		h.callbackAnswer(callback, err.Error())
		h.bot.Send(msg)
		return
	}

	if len(products) == 0 {
		msg = tgbotapi.NewMessage(callback.Message.Chat.ID, "товаров нет")
		h.callbackAnswer(callback, "товаров нет")
		h.bot.Send(msg)
		return
	}

	totalPages := (total + itemsPerPage - 1) / itemsPerPage

	var text strings.Builder
	for i, product := range products {
		fmt.Fprintf(&text, "%d. %s\nЦена: %.2f руб.\n%s\nВ наличии: %d\n\n",
			i+1, product.Name, product.Price, product.Description, product.Quantity)
	}
	fmt.Fprintf(&text, "сраница %d из %d\n", page+1, totalPages)

	editMsg := tgbotapi.NewEditMessageTextAndMarkup(callback.Message.Chat.ID, callback.Message.MessageID, text.String(), keyboard.CatalogMenu(page, totalPages, products))

	h.bot.Send(editMsg)
	h.callbackAnswer(callback, "done")
}

func (h *Handler) product_description(callback *tgbotapi.CallbackQuery) {
	id, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "product_description_"))
	Id := int64(id)

	product, err := h.storage.GetProduct(Id)
	if err != nil {
		h.callbackAnswer(callback, "товара не существует")
		return
	}

	text := fmt.Sprintf("%s\nЦена: %.2f руб.\n%s\nВ наличии: %d\n\n",
		product.Name, product.Price, product.Description, product.Quantity)

	if product.ImageUrl == "" {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, text)
		msg.ReplyMarkup = keyboard.ProductCartAddKeyboard(Id)
		h.bot.Send(msg)
		h.callbackAnswer(callback, "done")

	} else {
		path := "." + product.ImageUrl
		msg := tgbotapi.NewPhoto(callback.Message.Chat.ID, tgbotapi.FilePath(path))
		msg.Caption = text
		msg.ReplyMarkup = keyboard.ProductCartAddKeyboard(Id)
		h.bot.Send(msg)
		h.callbackAnswer(callback, "done")
	}

}

func (h *Handler) cart_add(callback *tgbotapi.CallbackQuery) {
	id, _ := strconv.ParseInt(strings.TrimPrefix(callback.Data, "cart_add_"), 10, 64)
	err := h.storage.AddToCart(callback.From.ID, id, 1)
	if err != nil {
		h.callbackAnswer(callback, err.Error())
	} else {
		h.callbackAnswer(callback, "done")
	}
}

func (h *Handler) cart_menu_order(callback *tgbotapi.CallbackQuery) {
	err := h.storage.Order(callback.From.ID)
	if err != nil {
		h.callbackAnswer(callback, err.Error())
	} else {
		editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "заказ оформлен")
		h.bot.Send(editMsg)
		h.callbackAnswer(callback, "done")
	}
}

func (h *Handler) cart_menu_clear_all(callback *tgbotapi.CallbackQuery) {
	h.storage.ClearCart(callback.From.ID)
	editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "корзина пуста")
	h.bot.Send(editMsg)
	h.callbackAnswer(callback, "done")
}
