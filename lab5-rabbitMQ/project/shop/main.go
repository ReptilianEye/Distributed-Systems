package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	. "example.com/hike-preparation/utils"
	rmq "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "orders_exchange"

func main() {
	conn, err := rmq.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // topic
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	FailOnError(err, "Failed to declare an exchange")

	var shopName string
	if len(os.Args) > 1 {
		shopName = os.Args[1]
	} else {
		shopName = "shop" + fmt.Sprintf("%d", time.Now().Nanosecond())
	}
	fmt.Println("Shop name:", shopName)

	go ListenForAdminMessages(ch, []string{"shop", "all"})

	var supportedProducts []string
	if len(os.Args) > 2 {
		supportedProducts = os.Args[2:]
	} else {
		var input string
		fmt.Println("Enter supported products separated by commas:")
		fmt.Scanf("%s", &input)
		supportedProducts = strings.Split(input, ",")
	}
	var queues = make([]rmq.Queue, len(supportedProducts))
	for i, product := range supportedProducts {
		queues[i], err = ch.QueueDeclare(
			product, // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		FailOnError(err, "Failed to declare a queue")
		fmt.Printf("Queue %s declared\n", product)
		ch.QueueBind(
			product,
			fmt.Sprintf("order.%s", product),
			exchangeName,
			false,
			nil,
		)
		fmt.Printf(
			"Queue %s bound to exchange %s with routing key %s\n",
			product,
			exchangeName,
			product,
		)
	}
	ch.Qos(
		1,
		0,
		true,
	)

	orders := consumeFromAnyQueue(ch, queues)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for o := range orders {
			var orderData Order
			err := json.Unmarshal(o.Body, &orderData)
			FailOnError(err, "Failed to unmarshal order data")
			fmt.Printf("Received order from %s on: %s\n", orderData.Hiker, orderData.Product)
			response := map[string]any{
				"hiker":   orderData.Hiker,
				"product": orderData.Product,
				"status":  "success",
				"shop":    shopName,
			}
			b, err := json.Marshal(response)
			FailOnError(err, "Failed to marshal response")
			err = ch.Publish(
				// "",
				exchangeName,
				o.ReplyTo,
				false,
				false,
				rmq.Publishing{
					ContentType:   "application/json",
					CorrelationId: o.CorrelationId,
					Body:          b,
				},
			)
			FailOnError(err, "Failed to send response message")
			o.Ack(false)
		}
		wg.Done()
	}()

	log.Println("Waiting for orders. To exit press CTRL+C")
	wg.Wait()
}

func consumeFromAnyQueue(ch *rmq.Channel, queues []rmq.Queue) <-chan rmq.Delivery {
	out := make(chan rmq.Delivery)
	for _, q := range queues {
		orders, err := ch.Consume(
			q.Name,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		FailOnError(err, "Failed to register a consumer")
		go func() {
			for o := range orders {
				out <- o
			}
		}()
	}
	return out
}
