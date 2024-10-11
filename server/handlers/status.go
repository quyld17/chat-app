package handlers

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/entities/status"
	"github.com/quyld17/chat-app/middlewares"
)

func UpdateStatus(c echo.Context, dbMySQL *sql.DB) error {
	ws, err := middlewares.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer ws.Close()

	userId := c.Get("user_id").(int)

	err = status.Update(dbMySQL, userId)
	if err != nil {
		log.Printf("Failed to update user status: %v", err)
		middlewares.SendWebSocketError(ws, "Failed to update user status")
		return nil
	}

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket error:", err)

			err = status.Remove(dbMySQL, userId)
			if err != nil {
				log.Println("Error removing user from status table:", err)
			}
			log.Printf("User %d is offline and removed from the status table.", userId)
			break
		}
	}

	return nil
}
