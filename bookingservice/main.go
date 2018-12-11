package bookingservice

import (
	"cloud-microservice-go/bookingservice/listener"
	"cloud-microservice-go/bookingservice/rest"
	"cloud-microservice-go/lib/configuration"
	"cloud-microservice-go/lib/msgqueue"
	msgqueue_amqp "cloud-microservice-go/lib/msgqueue/amqp"
	"cloud-microservice-go/lib/persistence/dblayer"
	"flag"
	"github.com/streadway/amqp"
)

func panicIfErr(err error)  {
	if err != nil {
		panic(err)
	}
}

func main() {

	var eventListener msgqueue.EventListener
	var eventEmitter msgqueue.EventEmitter

	confPath := flag.String("config", "./configuration/config.json", "path to config file")
	flag.Parse()
	
	config, _ := configuration.ExtractConfiguration(*confPath)

	switch config.MessageBrokerType {
	case "amqp":
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		panicIfErr(err)

		eventListener, err = msgqueue_amqp.NewAMQPEventListener(conn, "events", "bookings")
		panicIfErr(err)

		eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
		panicIfErr(err)
	default:
		panic("Bad message broker type:" + config.MessageBrokerType)
	}

	dbhandler, _ := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)
	processor := listener.EventProcessor{eventListener, dbhandler}
	processor.ProcessEvents()

	rest.ServeAPI(config.RestfulEndPoint, dbhandler, eventEmitter)
}
