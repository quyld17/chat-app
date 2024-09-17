package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/entities/messages"
	"github.com/quyld17/chat-app/entities/rooms"
	"github.com/quyld17/chat-app/entities/users"
	"github.com/quyld17/chat-app/middlewares"
)

type IncomingMessage struct {
	Message string `json:"message"`
}

var userConnections = make(map[int][]*websocket.Conn)

func Chat(c echo.Context, db *sql.DB) error {
	ws, err := middlewares.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer ws.Close()

	receiverIdStr := c.QueryParam("receiver_id")
	if receiverIdStr == "" {
		log.Printf("Receiver ID is missing")
		middlewares.SendWebSocketError(ws, "Receiver ID is required")
		return nil
	}

	receiverId, err := strconv.Atoi(receiverIdStr)
	if err != nil {
		log.Printf("Invalid Receiver ID: %v", err)
		middlewares.SendWebSocketError(ws, "Invalid Receiver ID")
		return nil
	}

	senderId, err := users.GetId(c, db)
	if err != nil {
		log.Printf("Failed to retrieve sender ID: %v", err)
		middlewares.SendWebSocketError(ws, "Failed to retrieve user ID")
		return nil
	}

	roomId, err := rooms.GetId(db, receiverId, senderId)
	if err != nil {
		log.Printf("Error retrieving or creating room: %v", err)
		middlewares.SendWebSocketError(ws, "Error retrieving or creating room")
		return nil
	}

	messageHistory, err := messages.GetChatHistory(db, roomId, 0)
	if err != nil {
		log.Printf("Error retrieving message history: %v", err)
		middlewares.SendWebSocketError(ws, "Error retrieving message history")
		return nil
	}

	if err := ws.WriteJSON(messageHistory); err != nil {
		log.Printf("Failed to send response via WebSocket: %v", err)
		middlewares.SendWebSocketError(ws, "Failed to send response")
		return nil
	}

	userConnections[senderId] = append(userConnections[senderId], ws)

	for {
		var incomingMsg IncomingMessage
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			middlewares.SendWebSocketError(ws, "Error reading message")
			break
		}

		if err := json.Unmarshal(msg, &incomingMsg); err != nil {
			log.Printf("Error parsing message: %v", err)
			middlewares.SendWebSocketError(ws, "Invalid message format")
			continue
		}

		err = messages.Save(db, roomId, senderId, incomingMsg.Message)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			middlewares.SendWebSocketError(ws, "Error saving message")
			continue
		}

		updatedMessages, err := messages.GetChatHistory(db, roomId, 0)
		if err != nil {
			log.Printf("Error retrieving updated message history: %v", err)
			middlewares.SendWebSocketError(ws, "Error retrieving updated message history")
			continue
		}

		for _, senderWs := range userConnections[senderId] {
			senderWs.WriteJSON(updatedMessages)
		}
		for _, receiverWs := range userConnections[receiverId] {
			receiverWs.WriteJSON(updatedMessages)
		}
	}

	return nil
}
