package builder

import (
	"cloud-microservice-go/lib/msgqueue"
	"cloud-microservice-go/lib/msgqueue/amqp"
	"log"
	"os"
)

func NewEventListenerFromEnvironment() (msgqueue.EventListener, error)  {
	var listener msgqueue.EventListener
	var err error

	if url := os.Getenv("AMQP_URL"); url != "" {
		log.Printf("connecting AMQP broker ar %s", url)
		listener, err = amqp.NewAMQPEventListenerFromEnvironment()
		if err != nil {
			return nil, err
		}
	}

	return listener, nil
}
