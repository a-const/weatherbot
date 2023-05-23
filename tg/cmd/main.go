package main

import (
	"fmt"
	"log"
	"tg/env"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("..")
	viper.SetConfigName("variables")
	return viper.ReadInConfig()
}

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Something wrong with config: %s", err.Error())
	}
	bot, err := tgbotapi.NewBotAPI(viper.GetString("apikey"))
	if err != nil {
		panic(err)
	}
	postgres := new(env.PostgresDB)
	if err := postgres.NewPostgreDB(env.Config{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.sslmode"),
	}); err != nil {
		panic(err)
	}
	weather := env.ReceiveWeatherData()
	users, _ := postgres.GetUsers()
	for _, user := range *users {
		fmt.Printf("user1: %d, %s, %t", user.ID, user.Username, user.Notifications)
	}

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		env.ProcessIncomeMessage(update, bot, postgres, weather)
		env.ProcessInlineMessage(update, bot, postgres, weather)
	}
}
