package main

import (
	"cloud-microservice-go/eventservice/rest"
	"cloud-microservice-go/lib/configuration"
	msgqueue_amqp "cloud-microservice-go/lib/msgqueue/amqp"
	"cloud-microservice-go/lib/persistence/dblayer"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main()  {
	configPath := flag.String("conf", `.\configuration\config.json`,"flag to set the path to the configuration json file")
	flag.Parse()
	// extract configuration
	config, _ := configuration.ExtractConfiguration(*configPath)
	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}

	emitter ,err := msgqueue_amqp.NewAMQPEventEmitter(conn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to database...")
	dbHandler, _ := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)
	httpError, httptlsError := rest.ServeAPI(config.RestfulEndPoint, config.RestfulTLSEndPint, dbHandler, emitter)
	fmt.Println("serve event service at http port: " + config.RestfulEndPoint)

	select {
	case err := <- httpError:
		log.Fatal("HTTP Error:", err)
	case err := <- httptlsError:
		log.Fatal("HTTPS Error:", err)
	}
}
