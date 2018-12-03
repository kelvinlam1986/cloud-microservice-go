package main

import (
	"fmt"
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

	_, err = channel.QueueDeclare("my_queue", true, false, false, false, nil)
	if err != nil {
		panic("error while declaring the queue: " + err.Error())
	}

	err = channel.QueueBind("my_queue", "#", "events", false, nil)
	if err != nil {
		panic("error while binding the queue: " + err.Error())
	}

	msgs, err := channel.Consume("my_queue", "", false, false, false, false, nil)
	if err != nil {
		panic("error while consuming the queue: " + err.Error())
	}

	for msg := range msgs {
		fmt.Println("message received: " + string(msg.Body))
		msg.Ack(false)
	}

}
