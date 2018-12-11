package amqp

import (
	"cloud-microservice-go/lib/msgqueue"
	"fmt"
	"github.com/streadway/amqp"
)

const eventNameHeader = "x-event-name"

type amqbEventListener struct {
	connection *amqp.Connection
	exchange string
	queue string
	mapper msgqueue.EventMapper
}

func (a *amqbEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	err = channel.ExchangeDeclare(a.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return  err
	}
	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not declare the queue %s: %s", a.queue, err)
	}

	return nil
}

func (l *amqbEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := l.connection.Channel()
	if err != nil {
		return nil, nil, err
	}

	for _, event := range eventNames {
		if err := channel.QueueBind(l.queue, event, l.exchange, false, nil); err != nil {
			return nil, nil, fmt.Errorf("could not bind the event %s to queue %s: %s", event, l.queue, err)
		}
	}

	msgs, err := channel.Consume(l.queue, "", false, false, false, false, nil)
	if err != nil {
		return  nil, nil, err
	}

	events := make(chan msgqueue.Event)
	errors := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers[eventNameHeader]
			if !ok {
				errors <- fmt.Errorf("msg did not contain %s header", eventNameHeader)
				msg.Nack(false, false)
				continue
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("header %s did not contain string", eventNameHeader)
				msg.Nack(false, false)
				continue
			}

			event, err := l.mapper.MapEvent(eventName, msg.Body)
			if err != nil {
				errors <- fmt.Errorf("could not unmarshall event %s, %s", eventName, err)
				msg.Nack(false, false)
				continue
			}

			events <- event
			msg.Ack(false)
		}
	}()

	return events, errors, nil
}

func NewAMQPEventListener(conn *amqp.Connection, exchange string, queue string) (msgqueue.EventListener, error) {
	listener := amqbEventListener{
		connection: conn,
		exchange: exchange,
		queue: queue,
		mapper: msgqueue.NewEventMapper(),
	}

	err := listener.setup()
	if err != nil {
		return nil, err
	}

	return &listener, nil
}

func (l *amqbEventListener) Mapper() msgqueue.EventMapper {
	return l.mapper
}