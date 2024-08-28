package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []byte)            // Broadcast channel

func handleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return err
	}
	defer ws.Close()

	// Register the new client
	clients[ws] = true

	// Handle incoming messages
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			delete(clients, ws)
			break
		}
		log.Printf("Received: %s\n", msg)

		// Send the message to the broadcast channel
		broadcast <- msg
	}

	return nil
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send the message to every client connected
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	e := echo.New()

	// WebSocket endpoint
	e.GET("/ws", handleWebSocket)

	// Start handling messages
	go handleMessages()

	e.Logger.Fatal(e.Start(":8080"))
}
