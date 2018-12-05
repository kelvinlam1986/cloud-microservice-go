package bookingservice

import (
	"cloud-microservice-go/bookingservice/listener"
	"cloud-microservice-go/lib/configuration"
	msgqueue_amqp "cloud-microservice-go/lib/msgqueue/amqp"
	"cloud-microservice-go/lib/persistence/dblayer"
	"flag"
	"github.com/streadway/amqp"
)

func main() {
	confPath := flag.String("config", "./configuration/config.json", "path to config file")
	flag.Parse()
	config, _ := configuration.ExtractConfiguration(*confPath)

	dbhandler, err := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)
	if err != nil {
		panic(err)
	}

	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}

	eventListener, err := msgqueue_amqp.NewAMQPEventListener(conn, "events")
	if err != nil {
		panic(err)
	}

	processor := &listener.EventProcessor{eventListener, dbhandler}
	processor.ProcessEvents()
}
