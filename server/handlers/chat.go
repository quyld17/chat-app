package handlers

// import (
// 	"database/sql"
// 	// "log"
// 	"net/http"
// 	// "strconv"

// 	"github.com/gorilla/websocket"
// 	"github.com/labstack/echo/v4"
// )

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func ChatWebSocket(c echo.Context, db *sql.DB) error {
	// conversationID := c.Param("conversation_id")
	// convID, err := strconv.Atoi(conversationID)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid conversation ID"})
	// }

	// ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	// if err != nil {
	// 	return err
	// }
	// defer ws.Close()

	// Handle WebSocket communication
	// for {
	// 	_, msg, err := ws.ReadMessage()
	// 	if err != nil {
	// 		log.Println("Error reading message:", err)
	// 		break
	// 	}

	// 	// Parse the message and save it to the database
	// 	var senderID int
	// 	var messageText string
	// 	// You can add logic to parse the `msg` and extract `senderID` and `messageText`

	// 	_, err = db.Exec(
	// 		`INSERT INTO messages (id, sender_id, message) 
	// 				VALUES (?, ?, ?);`,
	// 		convID, senderID, messageText,
	// 	)
	// 	if err != nil {
	// 		log.Println("Error saving message to DB:", err)
	// 		continue
	// 	}

	// 	// Broadcast the message to other connected clients (implement client tracking)
	// }

// 	return nil
// }
