package services

import (
	"weather/services/accuweather"

	"weather/services/weatherapi"
	"weather/services/yandex"
)

var numofservices float64 = 3

type WeatherToSend struct {
	Temp_max float64
	Temp_min float64
	Wind     float64
	Rain     int
	Pressure float64
}

func CalculateWeather(
	ya *yandex.YandexWeather,
	accu *accuweather.AccuweatherWeather,
	api_w *weatherapi.WeatherApiWeather,
) *WeatherToSend {

	weather := &WeatherToSend{}
	// Max temperature as average temp from all services
	weather.Temp_max = (float64(ya.Forecast.Parts[0].Temp_max) +
		accu.DailyForecasts[0].Temperature.Maximum.Value +
		api_w.Forecast.Forecastday[0].Day.Maxtemp_c) / numofservices
	// Min temperature as average temp from all services
	weather.Temp_min = (float64(ya.Forecast.Parts[0].Temp_min) +
		accu.DailyForecasts[0].Temperature.Minimum.Value +
		api_w.Forecast.Forecastday[0].Day.Mintemp_c) / numofservices
	// Pressure comes only from yandex
	weather.Pressure = ya.Forecast.Parts[0].Pressure_mm
	// ApiWeather gives wind speed only in km\h so we should calculate m\s by dividing by 3.6
	weather.Wind = ya.Forecast.Parts[0].Wind_speed + (api_w.Forecast.Forecastday[0].Day.Maxwind_kph / 3.6)
	//
	prec := (float64(ya.Forecast.Parts[0].Prec_mm) + api_w.Forecast.Forecastday[0].Day.Totalprecip_mm) / 2

	if prec < 5 {
		weather.Rain = 1
	} else if prec >= 5 && prec <= 15 {
		weather.Rain = 2
	} else if prec > 15 {
		weather.Rain = 3
	}
	return weather
}
