package handlers

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/entities/users"
	"github.com/quyld17/chat-app/middlewares"
)

func GetOnlineList(c echo.Context, db *sql.DB) error {
	ws, err := middlewares.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	list, err := users.GetOnlineList(c, db)
	if err != nil {
		log.Println("Error retrieving online user list:", err)
		return ws.WriteJSON(map[string]string{"error": "Unable to retrieve online users"})
	}

	err = ws.WriteJSON(list)
	if err != nil {
		log.Println("Error sending online users list via WebSocket:", err)
		return err
	}

	return nil
}
