package env

import (
	"fmt"
	"time"

	"github.com/enescakir/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func sendWeatherToAll(db *PostgresDB) error {
	weather := ReceiveWeatherData()
	//msg := buildWeatherMessage(,weather)
	return nil
}

func buildWeatherMessageByUsername(username string) {

	cur_time_string := time.Now().Format("02/01/2006")

	var rain_str string
	switch weather.Rain {
	case 1:
		rain_str = "Маловероятны"
	case 2:
		rain_str = "Возможны"
	case 3:
		rain_str = "Ожидаются"
	default:
		rain_str = " "
	}
	// msg := tgbotapi.NewMessage(
	msg_weather := tgbotapi.NewMessage(
		update.CallbackQuery.From.ID,
		fmt.Sprintf(viper.GetString("weather"),
			cur_time_string,
			emoji.RedTrianglePointedUp,
			weather.Temp_max,
			emoji.RedTrianglePointedDown,
			weather.Temp_min,
			emoji.LeafFlutteringInWind,
			weather.Wind,
			emoji.Thermometer,
			weather.Pressure,
			emoji.CloudWithRain,
			rain_str))
}

func buildWeatherMessagByChatId(update tgbotapi.Update, weather *Weather) tgbotapi.MessageConfig {
	cur_time_string := time.Now().Format("02/01/2006")

	var rain_str string
	switch weather.Rain {
	case 1:
		rain_str = "Маловероятны"
	case 2:
		rain_str = "Возможны"
	case 3:
		rain_str = "Ожидаются"
	default:
		rain_str = " "
	}
	msg_weather := tgbotapi.NewMessage(
		update.CallbackQuery.From.ID,
		fmt.Sprintf(viper.GetString("weather"),
			cur_time_string,
			emoji.RedTrianglePointedUp,
			weather.Temp_max,
			emoji.RedTrianglePointedDown,
			weather.Temp_min,
			emoji.LeafFlutteringInWind,
			weather.Wind,
			emoji.Thermometer,
			weather.Pressure,
			emoji.CloudWithRain,
			rain_str))
	return msg_weather
}

func ProcessIncomeMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *PostgresDB, weather *Weather) {
	if update.Message == nil {
		return
	}
	switch update.Message.Text {
	case "/start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("greetings"), emoji.WavingHand, emoji.Sun))
		msg.ReplyMarkup = GreetingsInline
		db.CreateUser(update.SentFrom().UserName)
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	case "/offnotifications":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(viper.GetString("offnotifications"), emoji.RightArrowCurvingDown))
		msg.ReplyMarkup = OffNotificationsInline
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}

func ProcessInlineMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *PostgresDB, weather *Weather) {
	if update.CallbackQuery == nil {
		return
	}
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	switch callback.Text {
	case "offnotifications":
		db.EditNotofications(update.SentFrom().UserName, false)
		msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Уведомления выключены!")
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	case "turnon":
		db.EditNotofications(update.SentFrom().UserName, true)
		msg_on := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Уведомления включены!")
		if _, err := bot.Send(msg_on); err != nil {
			panic(err)
		}
		_, err := bot.Send(buildWeatherMessage(update, weather))
		if err != nil {
			panic(err)
		}
		msgInfo := tgbotapi.NewMessage(update.CallbackQuery.From.ID, fmt.Sprintf(viper.GetString("notifinfo"), emoji.AlarmClock))
		_, err = bot.Send(msgInfo)
		if err != nil {
			panic(err)
		}

	}
}
