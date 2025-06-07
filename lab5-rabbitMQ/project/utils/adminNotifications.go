package utils

import (
	"fmt"
	"sync"

	rmq "github.com/rabbitmq/amqp091-go"
)

func ListenForAdminMessages(ch *rmq.Channel, keys []string) {
	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // auto-deleted
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare a queue")
	fmt.Printf("Queue %s declared\n", q.Name)

	for _, k := range keys {
		ch.QueueBind(
			q.Name,
			fmt.Sprintf("admin.%s", k),
			ExchangeName,
			false,
			nil,
		)
		FailOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to register a consumer")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for d := range msgs {
			fmt.Printf("Received following message from admin: %s\n", d.Body)
		}
		wg.Done()
	}()

	wg.Wait()
}
