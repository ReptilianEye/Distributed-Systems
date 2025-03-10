package main

import (
	"fmt"
	"log"
	"net"
)

const buffSize = 4096
const port = 8080

func main() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	udpAddr := net.UDPAddr{
		Port: port,
		IP:   nil,
	}

	var clients = make(map[string]net.Conn)
	udpListenConn, err := net.ListenUDP("udp", &udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpListenConn.Close()
	go handleUDPConnection(udpListenConn, clients)

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

		// Authenticate client
		nicknameByte := make([]byte, buffSize)
		_, err = tcpConn.Read(nicknameByte)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nickname := string(nicknameByte)
		clients[nickname] = tcpConn
		fmt.Printf("New client: %s\n", nickname)

		go handleTCPConnection(tcpConn, nickname, clients)
	}
}

func handleTCPConnection(conn net.Conn, clientNick string, otherClients map[string]net.Conn) {
	for {
		message := make([]byte, buffSize)
		_, err := conn.Read(message)
		if err != nil {
			fmt.Printf("Connection closed by client %s\n", clientNick)
			delete(otherClients, clientNick)
			return
		}
		forwardToOthersTCP(clientNick, message, otherClients)
	}
}

func handleUDPConnection(conn *net.UDPConn, otherClients map[string]net.Conn) {
	for {
		message := make([]byte, buffSize)
		_, addr, err := conn.ReadFromUDP(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Received message from UDP %s: %s\n", addr, message)
		forwardToOthersUDP(addr, message, otherClients)
	}
}
func forwardToOthersTCP(sender string, message []byte, otherClients map[string]net.Conn) {
	for nick, conn := range otherClients {
		if nick != sender {
			mess := []byte(sender + ": " + string(message))
			_, err := conn.Write(mess)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func forwardToOthersUDP(senderAddr *net.UDPAddr, message []byte, otherClients map[string]net.Conn) {
	udpConn := func(tcpConn net.Conn) *net.UDPConn {
		udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
			Port: retrievePort(tcpConn) + 1,
			IP:   nil,
		})
		if err != nil {
			log.Fatal(err)
		}
		return udpConn
	}
	var sender string
	for nick, conn := range otherClients {
		if retrievePort(conn) == senderAddr.Port {
			sender = nick
			break
		}
	}
	for _, conn := range otherClients {
		udpConn := udpConn(conn)
		if retrievePort(conn) != senderAddr.Port {
			_, err := udpConn.Write([]byte(sender + ": " + string(message)))
			if err != nil {
				fmt.Println(err)
			}
		}
		udpConn.Close()
	}
}

func retrievePort(conn net.Conn) int {
	var port int
	remoteAddr := conn.RemoteAddr().String()
	_, err := fmt.Sscanf(remoteAddr, "127.0.0.1:%d", &port)
	if err != nil {
		fmt.Println("Failed to extract port from remote address:", err)
		return 0
	}
	return port
}
