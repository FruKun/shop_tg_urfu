package keyboard

import (
	"fmt"

	"github.com/FruKun/shop_tg_bot_urfu/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenu() tgbotapi.ReplyKeyboardMarkup {
	buttons := [][]tgbotapi.KeyboardButton{
		{tgbotapi.NewKeyboardButton("каталог")},
		{tgbotapi.NewKeyboardButton("корзина")},
		{tgbotapi.NewKeyboardButton("помощь")},
	}
	return tgbotapi.NewReplyKeyboard(buttons...)
}

func getProductsDescriptionKeyboard(products []models.Product) []tgbotapi.InlineKeyboardButton {
	var keyboard []tgbotapi.InlineKeyboardButton
	for i, product := range products {
		keyboard = append(keyboard, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("🔎 %d", i+1), fmt.Sprintf("product_description_%d", product.Id)))
	}
	return keyboard
}

func ProductCartAddKeyboard(id int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("добавить в корзину", fmt.Sprintf("cart_add_%d", id))))
}

func CatalogMenu(page int, totalPages int, products []models.Product) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	switch page {
	case 0:
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", fmt.Sprintf("catalog_page_%d", totalPages)),
			tgbotapi.NewInlineKeyboardButtonData("Вперед ➡️", fmt.Sprintf("catalog_page_%d", page+1)),
		})
		buttons = append(buttons, getProductsDescriptionKeyboard(products))
	case totalPages:
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", fmt.Sprintf("catalog_page_%d", page-1)),
			tgbotapi.NewInlineKeyboardButtonData("Вперед ➡️", fmt.Sprintf("catalog_page_%d", 0)),
		})
		buttons = append(buttons, getProductsDescriptionKeyboard(products))
	default:
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", fmt.Sprintf("catalog_page_%d", page-1)),
			tgbotapi.NewInlineKeyboardButtonData("Вперед ➡️", fmt.Sprintf("catalog_page_%d", page+1)),
		})
		buttons = append(buttons, getProductsDescriptionKeyboard(products))
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func CartMenu() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("заказать", "cart_menu_order"),
		tgbotapi.NewInlineKeyboardButtonData("очистить", "cart_menu_clear_all"),
	})
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func CartMenuClear() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("очистить все", "cart_menu_clear_all"),
	})
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
