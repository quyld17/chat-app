package handlers

import (
	"log"
	"message-delivery/middlewares"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var userConnections = make(map[int]*websocket.Conn)

func SetUpConnection(c echo.Context) error {
	ws, err := middlewares.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer ws.Close()

	senderId := c.Get("user_id").(int)
	userConnections[senderId] = ws

	defer func() {
		delete(userConnections, senderId)
		log.Printf("User %d disconnected", senderId)
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("User %d read error: %v", senderId, err)
			break
		}
		log.Printf("Received from %d: %s", senderId, msg)
	}

	return nil
}
