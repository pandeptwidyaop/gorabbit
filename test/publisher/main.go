package main

import "github.com/pandeptwidyaop/gorabbit"

func main() {

	conect := "amqp://guest:guest@localhost:5672/"

	rb, err := gorabbit.New(conect, "test", "amq.topic")

	if err != nil {
		panic(err)
	}

	err = rb.Connect()
	if err != nil {
		panic(err)
	}

	err = rb.Publish("test", "application/json", []byte("Hello World!"))

	if err != nil {
		panic(err)
	}

}
