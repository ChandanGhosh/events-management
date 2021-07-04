package mqp

import (
	"encoding/json"
	"fmt"

	"github.com/chandanghosh/events-management/contracts"
	"github.com/chandanghosh/events-management/contracts/lib/msgqueue"
	"github.com/streadway/amqp"
)

type amqpEventListener struct {
	connection *amqp.Connection
	queue      string
}

func (a *amqpEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)
	return err

}

func (a *amqpEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	defer channel.Close()

	for _, eventName := range eventNames {
		if err = channel.QueueBind(a.queue, eventName, "events", false, nil); err != nil {
			return nil, nil, err
		}
	}
	msgChan, err := channel.Consume(a.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	eventsChan := make(chan msgqueue.Event)
	errorsChan := make(chan error)

	go func() {
		for msg := range msgChan {

			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errorsChan <- fmt.Errorf("message did not contain %s header", "x-event-name")
				msg.Nack(false, true)
				continue
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errorsChan <- fmt.Errorf("x-event-name is not a string, but of type %t", rawEventName)
				msg.Nack(false, false)
				continue
			}

			var event msgqueue.Event
			switch eventName {
			case "event.created":
				event = new(contracts.EventCreatedEvent)
			default:
				errorsChan <- fmt.Errorf("event type %s is unknown", eventName)
				continue
			}

			err := json.Unmarshal(msg.Body, event)
			if err != nil {
				errorsChan <- err
				continue
			}
			eventsChan <- event
		}
	}()
	return eventsChan, errorsChan, nil
}

func NewAMQPEventListener(conn *amqp.Connection, queue string) (msgqueue.EventListener, error) {
	listener := &amqpEventListener{
		connection: conn,
		queue:      queue,
	}
	if err := listener.setup(); err != nil {
		return nil, err
	}
	return listener, nil
}
