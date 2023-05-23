package accuweather

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// reference:
// https://developer.accuweather.com/accuweather-forecast-api/apis
// cabinet
// https://developer.accuweather.com/user/me/apps
// max requests for free: 50 per day

type AccuweatherWeather struct {
	DailyForecasts []struct {
		Temperature struct {
			Minimum struct {
				Value float64 `json:"Value"`
			} `json:"Minimum"`
			Maximum struct {
				Value float64 `json:"Value"`
			} `json:"Maximum"`
		} `json:"Temperature"`
	} `json:"DailyForecasts"`
}

func Load() *AccuweatherWeather {
	//struct for results
	accu_weather := &AccuweatherWeather{}
	// build request
	req_text := "http://dataservice.accuweather.com/forecasts/v1/daily/1day/16890?apikey="
	req_text += viper.GetString("accuweather_key")
	req_text += "&metric=true"
	req, _ := http.NewRequest("GET", req_text, nil)
	// send req, receive response, and unpack it to struct
	client := &http.Client{}
	if resp, err := client.Do(req); err != nil {
		logrus.Printf("accuWeather GET err! %s", err.Error())
	} else {
		defer resp.Body.Close()
		//bodyBytes, _ := io.ReadAll(resp.Body)
		//bodyString := string(bodyBytes)
		//fmt.Print(bodyString)
		if err := json.NewDecoder(resp.Body).Decode(&accu_weather); err != nil {
			logrus.Printf("error while decoding AccuWeather response! %s", err.Error())
		}
	}
	return accu_weather
}
