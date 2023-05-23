package env

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Weather struct {
	Temp_max float64 `json:"Temp_max"`
	Temp_min float64 `json:"Temp_min"`
	Wind     float64 `json:"Wind"`
	Rain     int     `json:"Rain"`
	Pressure float64 `json:"Pressure"`
}

func ReceiveWeatherData() *Weather {
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	weather := &Weather{}

	go func() {
		for d := range msgs {
			if err := json.Unmarshal(d.Body, weather); err != nil {
				log.Fatalf("error while decoding! %s", err.Error())
			}
		}
	}()

	return weather
}
