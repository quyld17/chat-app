package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/entities/status"
	"github.com/quyld17/chat-app/entities/users"
	"github.com/quyld17/chat-app/middlewares"
)

func UpdateStatus(c echo.Context, db *sql.DB) error {
	ws, err := middlewares.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	userID, err := users.GetID(c, db)
	if err != nil {
		log.Println("Error retrieving user ID:", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	err = status.Update(db, userID)
	if err != nil {
		log.Println("Error updating status to online:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update status to online"})
	}

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
