# gorabbit
## Rabbit MQ Publish & Subscribe

Simple library for AMQP Rabbit MQ publish subscribe

## Instalation
Install gorabbit depedency on your projects

```sh
go get github.com/pandeptwidyaop/gorabbit
```

## Usage

The following samples will assist you to become as comfortable as possible with gorabbit library.

```go
import "github.com/pandeptwidyaop/gorabbit"
```

### Listen/Subscribe

```go
func main(){
    // create infinite chan for listen queue
    forever := make(chan bool)

    mq, err := gorabbit.New(
        "amqp://user:password@host:port",
        "myQueueName",
        "myExchangeName",
    )

    if err != nil {
        panic(err)
    }

    // start connection
    err = mq.Connect()

    if err != nil {
        panic(err)
    }

    // binding all routing key

    err = mq.Bind([]string{"routing_a","routing_b","routing_c"})

    if err != nil {
        panic(err)
    }

    deliveries, err := mq.Consume()

    if err != nil {
        panic(err)
    }

    log.Println("Waiting for messages")

    for q, d := range deliveries {
        go mq.HandleConsumedDeliveries(q, d, handleConsume)
    }

    <-forever
}

// Handling messages
func handleConsume(mq gorabbit.RabbitMQ, queue string, deliveries <-chan amqp.Delivery){
    for d := range deliveries {
        switch d.RoutingKey {
        case "routing_a": 
            log.Println("message come from routing_a")
            // IMPORTANT !!!
            // Acknowledge the message
            // this will be called when success proceed
            // the message, if were not, the message will be still queued
            ackMessage(d)
        case "routing_b":
            log.Println("message come from routing_b")
            ackMessage(d)
        case "routing_c":
            log.Println("message come from routing_c")
            ackMessage(d)
        }
    }
}

func ackMessage(msg *amqp.Delivery) {
    if err := d.Ack(); err != nil {
        log.Fatal("Failed to acknowledge !")
    }
}
```

### Publish Message

Create struct for message body
```go

type Message struct {
    Name string `json:"name"`
    Address string `json:"address"`
}

```
Publish message
```go
m := Message{
    Name : "pande",
    Address : "address"
}

jsonMessage, err := json.Marshal(m)

if err != nil {
    panic(err)
}

mq, err := gorabbit.New(
    "amqp://user:password@host:port",
    "myQueueName",
    "myExchangeName",
)

if err != nil {
    panic(err)
}

// start connection
err = mq.Connect()

if err != nil {
    panic(err)
}

err = mq.Publish("routing_1", "application/json", jsonMessage)

if err != nil {
    panic(err)
}
```