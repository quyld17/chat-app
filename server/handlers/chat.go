package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/entities/messages"
	"github.com/quyld17/chat-app/entities/rooms"
	"github.com/quyld17/chat-app/middlewares"
)

type IncomingMessage struct {
	ReceiverId int    `json:"receiver_id"`
	Message    string `json:"message"`
}

func GetChatHistory(c echo.Context, db *sql.DB) error {
	offset := 0
	limit := 20

	receiverIdStr := c.QueryParam("receiver_id")
	if receiverIdStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Receiver ID is missing")
	}
	receiverId, err := strconv.Atoi(receiverIdStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Receiver ID")
	}

	offsetStr := c.QueryParam("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid offset")
		}
	}

	senderId := c.Get("user_id").(int)

	roomId, err := rooms.GetId(db, receiverId, senderId)
	if err == sql.ErrNoRows {
		roomId, err = rooms.Create(db, receiverId, senderId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create room")
		}
	}

	chatHistory, err := messages.GetHistory(db, roomId, offset, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error retrieving chat history")
	}

	for i, j := 0, len(chatHistory)-1; i < j; i, j = i+1, j-1 {
		chatHistory[i], chatHistory[j] = chatHistory[j], chatHistory[i]
	}

	return c.JSON(http.StatusOK, chatHistory)
}

var userConnections = make(map[int]*websocket.Conn)

func Chat(c echo.Context, db *sql.DB) error {
	ws, err := middlewares.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer ws.Close()

	senderId := c.Get("user_id").(int)

	userConnections[senderId] = ws

	defer delete(userConnections, senderId)

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
		receiverId := incomingMsg.ReceiverId

		roomId, err := rooms.GetId(db, receiverId, senderId)
		if err != nil {
			log.Printf("Error retrieving or creating room: %v", err)
			middlewares.SendWebSocketError(ws, "Error retrieving room")
			return nil
		}

		err = messages.Save(db, roomId, senderId, incomingMsg.Message)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			middlewares.SendWebSocketError(ws, "Error saving message")
			continue
		}

		response := map[string]interface{}{
			"sender_id": senderId,
			"message":   incomingMsg.Message,
		}

		if receiverWs, ok := userConnections[receiverId]; ok {
			if err := receiverWs.WriteJSON(response); err != nil {
				log.Printf("Error sending updated messages to receiver: %v", err)
			}
		}
	}

	return nil
}
