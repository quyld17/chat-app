package middlewares

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendWebSocketError(ws *websocket.Conn, message string) {
	errMsg := map[string]string{"error": message}
	ws.WriteJSON(errMsg)
}
