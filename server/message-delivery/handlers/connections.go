package handlers

import (
	"log"
	"message-delivery/middlewares"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var userConnections = map[int]map[*websocket.Conn]struct{}{}

func SetUpConnection(c echo.Context) error {
	ws, err := middlewares.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer ws.Close()

	senderId := c.Get("user_id").(int)
	if _, ok := userConnections[senderId]; !ok {
		userConnections[senderId] = make(map[*websocket.Conn]struct{})
	}
	userConnections[senderId][ws] = struct{}{}

	defer func() {
		if conns, ok := userConnections[senderId]; ok {
			delete(conns, ws)

			if len(conns) == 0 {
				delete(userConnections, senderId)
			}
		}
	}()

	return nil
}
