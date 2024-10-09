package handlers

import (
	"context"
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

// var userConnections = make(map[int][]*websocket.Conn)

var ctx = context.Background()

func Chat(c echo.Context, db *sql.DB, dbRedis *redis.Client) error {
	ws, err := middlewares.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}
	defer ws.Close()

	senderId := c.Get("user_id").(int)

	// userConnections[senderId] = append(userConnections[senderId], ws)

	// defer delete(userConnections, senderId)

	err = dbRedis.SAdd(ctx, "user:"+strconv.Itoa(senderId)+":connections", ws.RemoteAddr().String()).Err()
	if err != nil {
		log.Printf("Error storing user connection in Redis: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error storing connection")
	}

	defer func() {
		// Remove the connection when the WebSocket closes
		dbRedis.SRem(ctx, "user:"+strconv.Itoa(senderId)+":connections", ws.RemoteAddr().String())
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
		roomId, err := rooms.GetId(db, receiverId, senderId)
		if err != nil {
			log.Printf("Error retrieving room ID: %v", err)
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

		// if connections, ok := userConnections[receiverId]; ok {
		// 	for _, conn := range connections {
		// 		if err := conn.WriteJSON(response); err != nil {
		// 			log.Printf("Error sending new message to receiver: %v", err)
		// 		}
		// 	}
		// }

		receiverConnections, err := dbRedis.SMembers(ctx, "user:"+strconv.Itoa(receiverId)+":connections").Result()
		if err != nil {
			log.Printf("Error retrieving connections from Redis: %v", err)
			return nil
		}

		// Iterate through each connection and send the message
		for _, connAddr := range receiverConnections {
			// Create a temporary WebSocket connection for each address
			// This is a simulation; normally, you'd need to manage these connections separately
			tempConn, _, err := websocket.DefaultDialer.Dial("ws://"+connAddr, nil)
			if err != nil {
				log.Printf("Error connecting to receiver: %v", err)
				continue
			}
			defer tempConn.Close()

			if err := tempConn.WriteJSON(response); err != nil {
				log.Printf("Error sending new message to receiver: %v", err)
			}
		}
	}

	return nil
}
