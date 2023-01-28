package gorabbit

import (
	"context"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Context *RabbitMQ

	ErrConRefused   = errors.New("connection refused by server")
	ErrNotConnected = errors.New("not connected to server")
)

type RabbitMQ struct {
	connection       *amqp.Connection
	channel          *amqp.Channel
	queue            amqp.Queue
	err              chan error
	connectionString string
	exchange         string
	bindings         []string
	name             string
}

// New instance
func New(connectionString string, name string, exchange string) (*RabbitMQ, error) {
	mq := RabbitMQ{
		err:              make(chan error),
		connectionString: connectionString,
		exchange:         exchange,
		name:             name,
	}

	Context = &mq

	return &mq, nil
}

func (mq *RabbitMQ) Connect() error {
	var err error
	log.Println("[gorabbit] connecting to server")
	mq.connection, err = amqp.Dial(mq.connectionString)

	if err != nil {
		return err
	}

	mq.channel, err = mq.connection.Channel()

	if err != nil {
		return err
	}

	go func() {
		<-mq.channel.NotifyClose(make(chan *amqp.Error))
		mq.err <- ErrConRefused
	}()

	return nil
}

func (mq *RabbitMQ) Bind(bindings []string) error {
	var err error
	if mq.channel == nil {
		return ErrNotConnected
	}

	mq.bindings = bindings

	mq.queue, err = mq.channel.QueueDeclare(
		mq.name,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	for _, bind := range mq.bindings {
		log.Printf("[gorabbit] binding %s to channel \n", bind)
		err = mq.channel.QueueBind(
			mq.queue.Name,
			bind,
			mq.exchange,
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (mq *RabbitMQ) Reconnect() error {
	log.Println("[gorabbit] reconnecting to server")
	var err error

	err = mq.Connect()
	if err != nil {
		return err
	}

	err = mq.Bind(mq.bindings)
	if err != nil {
		return err
	}

	return nil
}

func (mq *RabbitMQ) Consume() (map[string]<-chan amqp.Delivery, error) {
	m := make(map[string]<-chan amqp.Delivery)
	for _, q := range mq.bindings {
		deliveries, err := mq.channel.Consume(
			mq.queue.Name,
			"",
			false, //auto acknowledge
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}
		m[q] = deliveries
	}
	return m, nil
}

func (mq *RabbitMQ) HandleConsumedDeliveries(q string, delivery <-chan amqp.Delivery, fn func(RabbitMQ, string, <-chan amqp.Delivery)) {
	for {
		go fn(*mq, q, delivery)
		if err := <-mq.err; err != nil {
			// Hard fail during connection loss
			panic(err)
			// err := mq.Reconnect()
			// if err != nil {
			// 	log.Println(err)
			// 	panic(err)
			// }
			// deliveries, err := mq.Consume()
			// if err != nil {
			// 	panic(err)
			// }

			// delivery = deliveries[q]
		}
	}
}

func (mq *RabbitMQ) Publish(event string, contentType string, message []byte) error {
	select {
	case err := <-mq.err:
		if err != nil {
			err := mq.Reconnect()
			if err != nil {
				return err
			}
		}
	default:
	}

	m := amqp.Publishing{
		ContentType: contentType,
		Body:        message,
	}

	ctx := context.Background()

	return mq.channel.PublishWithContext(
		ctx,
		mq.exchange,
		event,
		false,
		false,
		m,
	)
}
