package main

import (
	"fmt"
	"net"
)

const buffSize = 4096

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	var clients = make(map[string]net.Conn)
	// Accept incoming connections and handle them
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(clients) >= 10 {
			fmt.Println("Too many clients")
			conn.Close()
			continue
		}

		// Authenticate client
		nicknameByte := make([]byte, 1024)
		_, err = conn.Read(nicknameByte)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nickname := string(nicknameByte)
		clients[nickname] = conn
		fmt.Printf("New client: %s\n", nickname)

		// Handle the connection in a new goroutine
		go handleConnection(conn, nickname, clients)
	}
}

func handleConnection(conn net.Conn, clientNick string, otherClients map[string]net.Conn) {
	// Close the connection when we're done
	defer conn.Close()
	for {
		// Read incoming data
		message := make([]byte, buffSize)
		_, err := conn.Read(message)
		if err != nil {
			fmt.Printf("Connection closed by client %s\n", clientNick)
			delete(otherClients, clientNick)
			return
		}
		sendToAllOther(clientNick, message, otherClients)
	}
}

func sendToAllOther(nickname string, message []byte, otherClients map[string]net.Conn) {
	for nick, conn := range otherClients {
		if nick != nickname {
			mess := []byte(nickname + ": " + string(message))
			_, err := conn.Write(mess)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
