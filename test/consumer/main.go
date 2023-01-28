package main

import (
	"log"

	"github.com/pandeptwidyaop/gorabbit"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	con := "amqp://guest:guest@localhost:5672/"

	rb, err := gorabbit.New(con, "test", "amq.topic")
	if err != nil {
		panic(err)
	}

	err = rb.Connect()

	if err != nil {
		panic(err)
	}

	err = rb.Bind([]string{"test"})

	if err != nil {
		panic(err)
	}

	msgs, err := rb.Consume()

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	for q, delivery := range msgs {
		go rb.HandleConsumedDeliveries(q, delivery, handler)
	}

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}

func handler(mq gorabbit.RabbitMQ, q string, delivery <-chan amqp.Delivery) {
	for d := range delivery {
		log.Printf("Received a message: %s", d.Body)
		err := d.Ack(false)
		if err != nil {
			panic(err)
		}
	}
}
