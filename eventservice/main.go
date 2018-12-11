package main

import (
	"cloud-microservice-go/eventservice/rest"
	"cloud-microservice-go/lib/configuration"
	"cloud-microservice-go/lib/msgqueue"
	msgqueue_amqp "cloud-microservice-go/lib/msgqueue/amqp"
	"cloud-microservice-go/lib/persistence/dblayer"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
)

func main()  {
	var eventEmitter msgqueue.EventEmitter

	configPath := flag.String("conf", `.\configuration\config.json`,"flag to set the path to the configuration json file")
	flag.Parse()
	// extract configuration
	config, _ := configuration.ExtractConfiguration(*configPath)

	switch config.MessageBrokerType {
	case "amqp":
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		if err != nil {
			panic(err)
		}

		eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
		if err != nil {
			panic(err)
		}
	default:
		panic("Bad message broker type: " + config.MessageBrokerType)
	}

	fmt.Println("Connecting to database...")
	dbHandler, _ := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)

	fmt.Println("Serving API...")
	httpError := rest.ServeAPI(config.RestfulEndPoint, dbHandler, eventEmitter)

	if httpError != nil {
		panic(httpError)
	}
}
