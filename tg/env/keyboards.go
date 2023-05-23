package env

import (
	"fmt"

	"github.com/enescakir/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var GreetingsInline = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Включить %v", emoji.CheckMarkButton), "turnon"),
	),
)

var OffNotificationsInline = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Выключить уведомления %v", emoji.CrossMark), "offnotifications"),
	),
)
