package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	. "example.com/hike-preparation/utils"
	rmq "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := rmq.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		ExchangeName, // name
		"topic",      // topic
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	FailOnError(err, "Failed to declare an exchange")

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
	ch.QueueBind(
		q.Name,
		"#",
		ExchangeName,
		false,
		nil,
	)

	go adminNotifications(ch)

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
			fmt.Printf("Received %s\n", d.Body)
		}
		wg.Done()
	}()

	wg.Wait()
}

func adminNotifications(ch *rmq.Channel) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(
			"Enter recipient (hiker, shop, all) and the message to send (e.g. 'hiker Hello hiker!'):",
		)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		fields := strings.Fields(input)
		fmt.Println(fields)
		if len(fields) < 2 {
			fmt.Println("Invalid input. Please try again.")
			continue
		}
		recipient := fields[0]
		message := strings.Join(fields[1:], " ")

		fmt.Println("Sending message to", recipient, "with message:", message)
		recipient = strings.TrimSpace(recipient)
		message = strings.TrimSpace(message)

		if !slices.Contains([]string{"hiker", "shop", "all"}, recipient) {
			fmt.Println("Invalid recipient. Please try again.")
			continue
		}

		ch.Publish(
			ExchangeName,
			fmt.Sprintf("admin.%s", recipient),
			false,
			false,
			rmq.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		)
	}
}
