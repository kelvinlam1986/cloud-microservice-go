package main

import (
	"github.com/streadway/amqp"
	"os"
)

func main()  {
	amqpUrl := os.Getenv("AMQP_URL")
	if amqpUrl == "" {
		amqpUrl = "amqp://guest:guest@localhost:5672"
	}
	connection, err := amqp.Dial(amqpUrl);
	if err != nil {
		panic("Can not established AMQP connection: " + err.Error())
	}

	defer connection.Close();

	channel, err := connection.Channel()
	if err != nil {
		panic("could not open channel: " + err.Error())
	}

	err = channel.ExchangeDeclare("events", "topic", true,
		false, false, false, nil)
	if err != nil {
		panic(err)
	}

	message := amqp.Publishing{
		Body: []byte("Hello World"),
	}

	err = channel.Publish("events", "some-routing-key", false, false, message)
	if err != nil {
		panic("error while publishing message: " + err.Error())
	}
}
