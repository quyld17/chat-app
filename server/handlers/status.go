package handlers

// import (
// 	"database/sql"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// 	"github.com/labstack/echo/v4"
// )

// var onlineUsers = make(map[string]*websocket.Conn)
// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func OnlineStatus(c echo.Context, db *sql.DB) error {
// 	// Upgrade HTTP to WebSocket
// 	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer ws.Close()

// 	// Retrieve token and validate user
// 	token := c.QueryParam("token")
// 	userID, err := validateToken(token)
// 	if err != nil {
// 		log.Println("Invalid token:", err)
// 		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
// 	}

// 	// Add user to the online users map
// 	onlineUsers[userID] = ws
// 	log.Printf("User %s is online.", userID)

// 	// Listen for any disconnection or errors
// 	for {
// 		_, _, err := ws.ReadMessage()
// 		if err != nil {
// 			log.Println("WebSocket error:", err)
// 			delete(onlineUsers, userID)
// 			log.Printf("User %s is offline.", userID)
// 			break
// 		}
// 	}

// 	return nil
// }

// func validateToken(token string) (string, error) {
// 	// Replace with actual token validation logic
// 	// For now, let's assume the token is the user ID itself
// 	return token, nil
// }
