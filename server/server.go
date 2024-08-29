package main

import (
	// "context"
	"net/http"
	"time"

	// "github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quyld17/chat-app/routers"
	"github.com/quyld17/chat-app/services/databases"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	dbMySQL := databases.NewMySQL()

	// // MongoDB Setup
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// client, err := mongo.Connect(context.Background(), clientOptions)
	// if err != nil {
	// 	panic(err)
	// }
	// dbMongo := client.Database("chatdb")

	// // Redis Setup
	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

	router := echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true,
		MaxAge:           int(24 * time.Hour.Seconds()),
	}))

	routers.RegisterAPIHandlers(router, dbMySQL)

	router.Logger.Fatal(router.Start(":8080"))
}

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// var clients = make(map[*websocket.Conn]bool) // Connected clients
// var broadcast = make(chan []byte)            // Broadcast channel

// func handleWebSocket(c echo.Context) error {
// 	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
// 	if err != nil {
// 		log.Println("Failed to upgrade connection:", err)
// 		return err
// 	}
// 	defer ws.Close()

// 	// Register the new client
// 	clients[ws] = true

// 	// Handle incoming messages
// 	for {
// 		_, msg, err := ws.ReadMessage()
// 		if err != nil {
// 			log.Printf("Error reading message: %v\n", err)
// 			delete(clients, ws)
// 			break
// 		}
// 		log.Printf("Received: %s\n", msg)

// 		// Send the message to the broadcast channel
// 		broadcast <- msg
// 	}

// 	return nil
// }

// func handleMessages() {
// 	for {
// 		// Grab the next message from the broadcast channel
// 		msg := <-broadcast

// 		// Send the message to every client connected
// 		for client := range clients {
// 			err := client.WriteMessage(websocket.TextMessage, msg)
// 			if err != nil {
// 				log.Printf("Error sending message: %v\n", err)
// 				client.Close()
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }

// func main() {
// 	e := echo.New()

// 	// WebSocket endpoint
// 	e.GET("/ws", handleWebSocket)

// 	// Start handling messages
// 	go handleMessages()

// 	e.Logger.Fatal(e.Start(":8080"))
// }
