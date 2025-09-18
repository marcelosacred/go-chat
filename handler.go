package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var clients = make(map[net.Conn]string)

var history []string
const historyLimit = 50

func handleConnection(conn net.Conn) {
	defer func() {
		fmt.Println("Client disconnected:", conn.RemoteAddr())
		delete(clients, conn)
		conn.Close()
	}()

	conn.Write([]byte("Welcome! Please enter your nickname: "))
	reader := bufio.NewScanner(conn)

	// Читаем никнейм
	reader.Scan()
	nickname := strings.TrimSpace(reader.Text())
	if nickname == "" {
		nickname = "Anonymous"
	}
	clients[conn] = nickname
	fmt.Println(nickname, "connected from", conn.RemoteAddr())
	conn.Write([]byte("Hello, " + nickname + "!\n"))

	// Отправляем историю сообщений новому клиенту
	for _, msg := range history {
		conn.Write([]byte(msg + "\n"))
	}

	// Читаем сообщения от клиента
	for reader.Scan() {
		msg := reader.Text()
		fullMsg := nickname + ": " + msg

		// Добавляем в историю
		history = append(history, fullMsg)
		if len(history) > historyLimit {
			history = history[1:]
		}

		// Рассылаем всем остальным
		broadcast(conn, fullMsg)
	}

	if err := reader.Err(); err != nil {
		fmt.Println("Read error:", err)
	}
}

// broadcast отправляет сообщение всем клиентам, кроме отправителя
func broadcast(sender net.Conn, msg string) {
	for client := range clients {
		if client != sender {
			_, _ = client.Write([]byte(msg + "\n"))
		}
	}
}
