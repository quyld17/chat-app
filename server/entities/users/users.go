package users

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"cratedAt"`
}

type Status struct {
	Id int `json:"id"`
}

func Authenticate(account Users, db *sql.DB) error {
	var storedPassword string

	err := db.QueryRow(`
        SELECT password 
        FROM users
        WHERE username = ?;
    `, account.Username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Invalid username! Please try again.")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(account.Password))
	if err != nil {
		return fmt.Errorf("Invalid username or password! Please try again.")
	}

	return nil
}

func Create(newUser Users, db *sql.DB) error {
	_, err := db.Exec(`
        INSERT INTO users (username, password) 
        VALUES (?, ?)
    `, newUser.Username, newUser.Password)
	if err != nil {
		return err
	}
	return nil
}

func GetID(c echo.Context, db *sql.DB) (int, error) {
	username := c.Get("username").(string)
	row := db.QueryRow(`
		SELECT id
		FROM users
		WHERE username = ?;
		`, username)
	var userID int
	if err := row.Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

func GetOnlineList(c echo.Context, db *sql.DB) ([]Users, error) {
	rows, err := db.Query(`
		SELECT 	status.user_id, 
				users.username
		FROM status
		JOIN users ON status.user_id = users.id;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []Users{}
	for rows.Next() {
		var user Users
		err := rows.Scan(&user.Id, &user.Username)
		if err != nil {
			return nil, err
		}
		list = append(list, user)
	}

	return list, nil
}
