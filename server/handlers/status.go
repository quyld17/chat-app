package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/entities/status"
	"github.com/quyld17/chat-app/entities/users"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func UpdateStatus(c echo.Context, db *sql.DB) error {
	// Upgrade the HTTP request to a WebSocket connection
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	userID, err := users.GetID(c, db)
	if err != nil {
		log.Println("Error retrieving user ID:", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	err = status.Update(db, userID, )
	if err != nil {
		log.Println("Error updating status to online:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update status to online"})
	}

	log.Printf("User %s is online.", userID)

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("WebSocket error:", err)

			// Remove the user's status from the database on disconnect
			err = status.Remove(db, userID)
			if err != nil {
				log.Println("Error removing user from status table:", err)
			}
			log.Printf("User %s is offline and removed from the status table.", userID)
			break
		}
	}

	return nil
}
