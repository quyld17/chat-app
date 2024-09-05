package middlewares

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	jwtService "github.com/quyld17/chat-app/services/jwt"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func JWTAuthorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := jwtService.GetToken(c)
		if tokenString == "" {
			tokenString = c.QueryParam("token")
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: Missing token")
			}
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: Invalid token")
		}

		username := jwtService.GetClaims(token, "username")
		if username == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid token claims: Missing username")
		}

		c.Set("username", username)

		return next(c)
	}
}

func SendWebSocketError(ws *websocket.Conn, message string) {
	errMsg := map[string]string{"error": message}
	ws.WriteJSON(errMsg)
}

// func Pagination(c echo.Context, itemsPerPage int) (int, error) {
// 	pageStr := c.QueryParam("page")
// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil {
// 		return 0, fmt.Errorf("Invalid page number")
// 	}

// 	offset := (page - 1) * itemsPerPage

// 	return offset, nil
// }
