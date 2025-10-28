package jwt

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetToken(c echo.Context) string {
	token := c.Request().Header.Get("Authorization")
	return token
}

func GetClaims(token *jwt.Token, key string) string {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}
	value, ok := claims[key].(string)
	if !ok {
		return ""
	}
	return value
}

func Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := GetToken(c)
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

		username := GetClaims(token, "name")
		if username == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid token claims: Missing username")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid token claims: Missing user_id")
		}
		userIdFloat, ok := claims["user_id"].(float64)
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid token claims: Invalid user_id format")
		}
		userId := int(userIdFloat)
		if !ok || userId == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid token claims: Missing user_id")
		}

		c.Set("username", username)
		c.Set("user_id", userId)

		return next(c)
	}
}
