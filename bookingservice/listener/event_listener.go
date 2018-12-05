package listener

import (
	"cloud-microservice-go/contracts"
	"cloud-microservice-go/lib/msgqueue"
	"cloud-microservice-go/lib/persistence"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type EventProcessor struct {
	EventListener msgqueue.EventListener
	Database persistence.DatabaseHandler
}

func (p *EventProcessor) ProcessEvents() error {
	log.Println("Listening to events...")
	received, errors, err := p.EventListener.Listen("event.created")
	if err != nil {
		return err
	}

	for {
		select {
		case evt := <-received:
			p.handleEvent(evt)
		case err := <-errors:
			log.Printf("received error while processing message: %s", err)
		}
	}
}

func (p *EventProcessor) handleEvent(event msgqueue.Event)  {
	switch e := event.(type) {
	case *contracts.EventCreatedEvent:
		log.Printf("event %s created: %s", e.ID, e)
		p.Database.AddEvent(persistence.Event{ID:bson.ObjectId(e.ID)})
	default:
		log.Printf("unknown event: %t", e)
	}
}