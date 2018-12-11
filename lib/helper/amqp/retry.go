package amqp

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func RetryConnect(amqpUrl string, retryInterval time.Duration) chan *amqp.Connection {
	result := make(chan *amqp.Connection)

	go func() {
		defer close(result)
		for {
			conn, err := amqp.Dial(amqpUrl)
			if err == nil {
				log.Println("connection successful established")
				result <- conn
				return
			}

			log.Printf("AMQP connection failed with error (retrying in %s): %s", retryInterval.String(), err)
			time.Sleep(retryInterval)
		}
	}()

	return result
}
