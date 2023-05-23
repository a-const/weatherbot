package yandex

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// reference
// https://yandex.ru/dev/weather/doc/dg/concepts/forecast-info.html#req-format
// cabinet:
// https://developer.tech.yandex.ru/services/18
// max requests for free: 50 per day

type YandexWeather struct {
	Forecast struct {
		Parts []struct {
			Temp_max    int     `json:"temp_max"`
			Temp_min    int     `json:"temp_min"`
			Wind_speed  float64 `json:"wind_speed"`
			Pressure_mm float64 `json:"pressure_mm"`
			Prec_mm     int     `json:"prec_mm"`
			Prec_period int     `json:"prec_period"`
		} `json:"parts"`
	} `json:"forecast"`
}

func Load() *YandexWeather {
	//struct for results
	ya_weather := &YandexWeather{}
	// build request
	req, _ := http.NewRequest("GET", "https://api.weather.yandex.ru/v2/informers?lat=40.177628&lon=44.512555", nil)
	req.Header.Set("X-Yandex-API-Key", viper.GetString("ya_key"))
	// send req, receive response, and unpack it to struct
	client := &http.Client{} // Make client instace in main function.
	if resp, err := client.Do(req); err != nil {
		logrus.Printf("Yandex GET err! %s", err.Error())
	} else if err := json.NewDecoder(resp.Body).Decode(&ya_weather); err != nil {
		logrus.Printf("Error while decoding yandex response! %s", err.Error())
	}
	return ya_weather
}
