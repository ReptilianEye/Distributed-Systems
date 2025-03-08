package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
)

const buffSize = 20000
const addr = "127.0.0.1"
const multicastAddr = "224.0.0.1:12345"

func main() {
	// Connect to the server
	tcpConn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	var tcpPort int
	_, err = fmt.Sscanf(tcpConn.LocalAddr().String(), "[::1]:%d", &tcpPort)
	if err != nil {
		fmt.Println("Failed to extract port from remote address:", err)
		tcpConn.Close()
		return
	}
	defer tcpConn.Close()

	udpConn, err := net.DialUDP(
		"udp",
		&net.UDPAddr{Port: tcpPort + 1, IP: net.ParseIP(addr)},
		&net.UDPAddr{Port: 8080, IP: net.ParseIP(addr)},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	go handleIncomingMessagesTCP(tcpConn)
	go handleIncomingMessagesUDP(tcpPort)
	// Authenticate
	nickname := fmt.Sprintf("client%d", rand.Intn(1000))
	fmt.Printf("Your nickname is: %s\n", nickname)
	_, err = tcpConn.Write([]byte(nickname))
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
		switch {
		case strings.HasPrefix(string(message), "U "):
			_, err = udpConn.Write(message)
		default:
			_, err = tcpConn.Write(message)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}
func handleIncomingMessagesTCP(conn net.Conn) {
	for {
		message := make([]byte, buffSize)
		_, err := conn.Read(message)
		if err != nil {
			fmt.Println(err)
			log.Fatal("Server closed the connection")
			return
		}
		fmt.Println("(TCP)", string(message))
	}
}

func handleIncomingMessagesUDP(port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(addr),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Listening for UDP messages on port", port)
	for {
		// Read incoming data
		message := make([]byte, buffSize)
		_, _, err := conn.ReadFromUDP(message)
		if err != nil {
			fmt.Println(err)
			log.Fatal("Server closed the connection")
			return
		}
		// Print the incoming data
		fmt.Println("(UDP)", string(message))
	}
}
