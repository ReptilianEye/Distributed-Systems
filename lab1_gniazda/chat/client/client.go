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
	defer tcpConn.Close()

	var tcpPort int
	_, err = fmt.Sscanf(tcpConn.LocalAddr().String(), "[::1]:%d", &tcpPort)
	if err != nil {
		fmt.Println("Failed to extract port from remote address:", err)
		tcpConn.Close()
		return
	}
	fmt.Printf("Client TCP runs on port %d\n", tcpPort)

	udpSend, err := net.DialUDP(
		"udp",
		&net.UDPAddr{Port: tcpPort, IP: nil},
		&net.UDPAddr{Port: 8080, IP: nil},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Client UDP runs on port", tcpPort+1)

	defer udpSend.Close()

	udpListen, err := net.ListenUDP("udp", &net.UDPAddr{Port: tcpPort + 1, IP: nil})
	if err != nil {
		log.Fatal(err)
	}
	defer udpListen.Close()

	multicastAddr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
	}

	multicastListen, err := net.ListenMulticastUDP(
		"udp",
		nil,
		multicastAddr,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer multicastListen.Close()

	go handleIncomingMessagesTCP(tcpConn)
	go handleIncomingMessagesUDP(udpListen)
	go handleIncomingMessagesUDP(multicastListen)
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
			_, err = udpSend.Write(message)
		case strings.HasPrefix(string(message), "M "):
			c, err := net.DialUDP(
				"udp",
				nil,
				multicastAddr,
			)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = c.Write(message)

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

func handleIncomingMessagesUDP(conn *net.UDPConn) {
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
