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
	"github.com/redis/go-redis/v9"
)

type IncomingMessage struct {
	ReceiverId int    `json:"receiver_id"`
	Message    string `json:"message"`
}

func GetChatHistory(c echo.Context, dbMySQL *sql.DB) error {
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

	roomId, err := rooms.GetId(dbMySQL, receiverId, senderId)
	if err == sql.ErrNoRows {
		roomId, err = rooms.Create(dbMySQL, receiverId, senderId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create room")
		}
	}

	chatHistory, err := messages.GetHistory(dbMySQL, roomId, offset, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error retrieving chat history")
	}

	for i, j := 0, len(chatHistory)-1; i < j; i, j = i+1, j-1 {
		chatHistory[i], chatHistory[j] = chatHistory[j], chatHistory[i]
	}

	return c.JSON(http.StatusOK, chatHistory)
}

var userConnections = make(map[int][]*websocket.Conn)

func Chat(c echo.Context, dbMySQL *sql.DB, dbRedis *redis.Client) error {
	ws, err := middlewares.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer ws.Close()

	senderId := c.Get("user_id").(int)

	userConnections[senderId] = append(userConnections[senderId], ws)

	defer func() {
		if connections, ok := userConnections[senderId]; ok {
			for i, conn := range connections {
				if conn == ws {
					userConnections[senderId] = append(connections[:i], connections[i+1:]...)
					break
				}
			}
			if len(userConnections[senderId]) == 0 {
				delete(userConnections, senderId)
			}
		}
	}()

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
		roomId, err := rooms.GetId(dbMySQL, receiverId, senderId)
		if err != nil {
			log.Printf("Error retrieving room ID: %v", err)
			middlewares.SendWebSocketError(ws, "Error retrieving room")
			return nil
		}

		err = messages.Save(dbMySQL, roomId, senderId, incomingMsg.Message)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			middlewares.SendWebSocketError(ws, "Error saving message")
			continue
		}

		response := map[string]interface{}{
			"sender_id": senderId,
			"message":   incomingMsg.Message,
		}

		if connections, ok := userConnections[receiverId]; ok {
			for _, conn := range connections {
				if err := conn.WriteJSON(response); err != nil {
					log.Printf("Error sending new message to receiver: %v", err)
				}
			}
		}
	}

	return nil
}
