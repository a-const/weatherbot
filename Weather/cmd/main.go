package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"weather/services"
	"weather/services/accuweather"
	"weather/services/weatherapi"
	"weather/services/yandex"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func initConfig() error {
	viper.AddConfigPath("..")
	viper.SetConfigName("keys")
	return viper.ReadInConfig()
}

func main() {
	fmt.Print("Weather parser 0.1\n\n")
	var ya_weather *yandex.YandexWeather
	var accu_weather *accuweather.AccuweatherWeather
	var api_weather *weatherapi.WeatherApiWeather

	if err := initConfig(); err != nil {
		log.Fatalf("Something wrong with config: %s", err.Error())
	}

	ya_weather = yandex.Load()
	accu_weather = accuweather.Load()
	api_weather = weatherapi.Load()

	if ya_weather != nil {
		fmt.Printf("Yandex response: \nTemperature max: %d\nTemperature min: %d\nWind speed: %f\nPressure mm: %f\nPrec mm: %d\nPrec period: %d\n\n",
			ya_weather.Forecast.Parts[0].Temp_max,
			ya_weather.Forecast.Parts[0].Temp_min,
			ya_weather.Forecast.Parts[0].Wind_speed,
			ya_weather.Forecast.Parts[0].Pressure_mm,
			ya_weather.Forecast.Parts[0].Prec_mm,
			ya_weather.Forecast.Parts[0].Prec_period)
	}
	if accu_weather != nil {
		fmt.Printf("Accuweather response: \nTemperature max: %f\nTemperature min: %f\n\n ",
			accu_weather.DailyForecasts[0].Temperature.Maximum.Value,
			accu_weather.DailyForecasts[0].Temperature.Minimum.Value)
	}
	if api_weather != nil {
		fmt.Printf("WeatherApi response: \nTemperature max: %f\nTemperature min: %f\nMaxwind: %f\nPrecip mm: %f\n\n ",
			api_weather.Forecast.Forecastday[0].Day.Maxtemp_c,
			api_weather.Forecast.Forecastday[0].Day.Mintemp_c,
			api_weather.Forecast.Forecastday[0].Day.Maxwind_kph,
			api_weather.Forecast.Forecastday[0].Day.Totalprecip_mm,
		)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//body, _ := json.Marshal(accu_weather)
	weathertosend := services.CalculateWeather(ya_weather, accu_weather, api_weather)
	body, _ := json.Marshal(*weathertosend)
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)

}
