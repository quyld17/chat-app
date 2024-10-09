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
	CreatedAt time.Time `json:"created_at"`
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
			return fmt.Errorf("Invalid username or password! Please try again.")
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
        VALUES (?, ?);
    `, newUser.Username, newUser.Password)
	if err != nil {
		return err
	}

	return nil
}

func GetIdByUsername(c echo.Context, db *sql.DB, username string) (int, error) {
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

func GetOnlineList(c echo.Context, db *sql.DB, userId int) ([]Users, error) {
	rows, err := db.Query(`
		SELECT 	status.user_id, 
				users.username
		FROM status
		JOIN users ON status.user_id = users.id
		WHERE status.user_id != ?;`,
		userId)
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

func CheckOrCreateGoogleAccount(c echo.Context, db *sql.DB, email string) error {
	var userID int
	err := db.QueryRow(`
		SELECT id
		FROM users
		WHERE username = ?;
	`, email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = db.Exec(`
				INSERT INTO users (username, password, is_google_account)
				VALUES (?, "0", 1);
			`, email)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return nil
}

func CheckIsGoogleAccount(c echo.Context, db *sql.DB, username string) (int, error) {
	var isGoogleAccount int
	err := db.QueryRow(`
		SELECT is_google_account
		FROM users
		WHERE username = ?;
	`, username).Scan(&isGoogleAccount)
	if err != nil {
		return 0, err
	}

	return isGoogleAccount, nil
}