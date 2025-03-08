package main

import (
	"fmt"
	"log"
	"net"
)

const buffSize = 4096
const port = 8080
const addr = "127.0.0.1"

type conns struct {
	tcp net.Conn
	udp *net.UDPConn
}

func main() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	udpAddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(addr),
	}

	var clients = make(map[string]conns)
	conn, err := net.ListenUDP("udp", &udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Accept incoming connections and handle them
	for {
		tcpConn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(clients) >= 10 {
			fmt.Println("Too many clients")
			tcpConn.Close()
			continue
		}
		remoteAddr := tcpConn.RemoteAddr().String()
		var remotePort int
		_, err = fmt.Sscanf(remoteAddr, "[::1]:%d", &remotePort)
		if err != nil {
			fmt.Println("Failed to extract port from remote address:", err)
			tcpConn.Close()
			continue
		}
		udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
			Port: remotePort,
			IP:   net.ParseIP(addr),
		})
		if err != nil {
			log.Fatal(err)
		}
		defer udpConn.Close()

		// Authenticate client
		nicknameByte := make([]byte, buffSize)
		_, err = tcpConn.Read(nicknameByte)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nickname := string(nicknameByte)
		clients[nickname] = conns{tcpConn, udpConn}
		fmt.Printf("New client: %s\n", nickname)

		// Handle the connection in a new goroutine
		go handleUDPConnection(conn, nickname, clients)
		go handleTCPConnection(tcpConn, nickname, clients)
	}
}

func handleTCPConnection(conn net.Conn, clientNick string, otherClients map[string]conns) {
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
		sendToAllOthersTCP(clientNick, message, otherClients)
	}
}

func handleUDPConnection(conn *net.UDPConn, clientNick string, otherClients map[string]conns) {
	for {
		message := make([]byte, buffSize)
		_, addr, err := conn.ReadFromUDP(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Received message from UDP %s: %s\n", addr, message)
		sendToAllOthersUDP(clientNick, message, otherClients)
	}
}
func sendToAllOthersTCP(nickname string, message []byte, otherClients map[string]conns) {
	for nick, conn := range otherClients {
		if nick != nickname {
			mess := []byte(nickname + ": " + string(message))
			_, err := conn.tcp.Write(mess)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func sendToAllOthersUDP(nickname string, message []byte, otherClients map[string]conns) {
	for nick, conn := range otherClients {
		if nick != nickname {
			mess := []byte(nickname + ": " + string(message))
			_, err := conn.udp.Write(mess)
			if err != nil {
				fmt.Println(err)
			}

		}
	}
}
