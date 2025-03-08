package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
)

const buffSize = 4096

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	go handleIncomingMessages(conn)
	// Authenticate
	nickname := fmt.Sprintf("client%d", rand.Intn(1000))
	fmt.Printf("Your nickname is: %s\n", nickname)
	_, err = conn.Write([]byte(nickname))
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		message, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			return
		}
		// Send some data to the server
		_, err = conn.Write(message)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Close the connection
	}
}
func handleIncomingMessages(conn net.Conn) {
	for {
		// Read incoming data
		message := make([]byte, buffSize)
		_, err := conn.Read(message)
		if err != nil {
			fmt.Println(err)
			log.Fatal("Server closed the connection")
			return
		}
		// Print the incoming data
		fmt.Printf("%s\n", message)
	}
}
