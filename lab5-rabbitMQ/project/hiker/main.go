package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	. "example.com/hike-preparation/utils"
	"github.com/google/uuid"
	rmq "github.com/rabbitmq/amqp091-go"
)

var pendingOrders = make(map[string]any)
var supportedProducts []string = []string{"oxygen", "bots", "backpack"}

const exchangeName = "orders_exchange"

func handleReadingResponses(msgs <-chan rmq.Delivery) {
	for d := range msgs {
		var response map[string]any
		err := json.Unmarshal(d.Body, &response)
		if pendingOrders[d.CorrelationId] == nil {
			fmt.Printf("Ignoring response for unknown correlation ID %s\n", d.CorrelationId)
			continue
		}
		FailOnError(err, "Failed to unmarshal response")

		fmt.Printf("Received response: %v\n", response)
		pendingOrders[d.CorrelationId] = nil
	}
	fmt.Println("All orders processed")
}

func main() {
	conn, err := rmq.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare an exchange")

	// We want to receive responses from the shop
	replyQ, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // auto-deleted
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare an exchange")
	msgs, err := ch.Consume(
		replyQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to register a consumer")

	go handleReadingResponses(msgs)
	go ListenForAdminMessages(ch, []string{"hiker", "all"})

	var hikerName string
	if len(os.Args) > 1 {
		hikerName = os.Args[1]
	} else {
		hikerName = "hiker" + fmt.Sprintf("%d", time.Now().Nanosecond())
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter orders separated by commas (or 'exit' to quit):")
		fmt.Println("Supported products:", supportedProducts)
		order, _ := reader.ReadString('\n')
		order = strings.TrimSpace(order)
		if order == "exit" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
		products := strings.Split(order, ",")
		for _, product := range products {
			product = strings.TrimSpace(product)
			if product == "" {
				continue
			}
			if !slices.Contains(supportedProducts, product) {
				fmt.Printf("Product %s not supported. Skipping...\n", product)
				continue
			}
			corrId := uuid.New().String()
			orderPayload := Order{
				Hiker:   hikerName,
				Product: product,
			}
			pendingOrders[corrId] = orderPayload
			b, err := json.Marshal(orderPayload)
			FailOnError(err, "Failed to marshal order payload")
			err = ch.Publish(
				exchangeName,                     // exchange
				fmt.Sprintf("order.%s", product), // routing key
				false,                            // mandatory
				false,                            // immediate
				rmq.Publishing{
					ContentType:   "application/json",
					CorrelationId: corrId,
					ReplyTo:       replyQ.Name,
					Body:          b,
				})
			FailOnError(err, "Failed to publish a message")
			fmt.Printf("Order %s sent\n", product)
		}
	}
}
