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

// func GetDetails(userID int, db *sql.DB) (*User, *Address, error) {
// 	row, err := db.Query(`
// 		SELECT
// 			email,
// 			full_name,
// 			phone_number,
// 			gender,
// 			date_of_birth
// 		FROM users
// 		WHERE user_id = ?;
// 		`, userID)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	var user User
// 	if row.Next() {
// 		err := row.Scan(&user.Email, &user.FullName, &user.PhoneNumber, &user.Gender, &user.DateOfBirth)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 	}
// 	user.DateOfBirthString = user.DateOfBirth.Format("2006-01-02")

// 	row, err = db.Query(`
// 		SELECT
// 			city,
// 			district,
// 			ward,
// 			street,
// 			house_number
// 		FROM addresses
// 		WHERE
// 			user_id = ? AND
// 			is_default = 1;
// 		`, userID)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer row.Close()

// 	var address Address
// 	if row.Next() {
// 		err := row.Scan(&address.City, &address.District, &address.Ward, &address.Street, &address.HouseNumber)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 	}

// 	return &user, &address, nil
// }

// func ChangePassword(userID int, password, newPassword string, c echo.Context, db *sql.DB) error {
// 	row, err := db.Query(`
// 		SELECT password
// 		FROM users
// 		WHERE
// 			user_id = ? AND
// 			password = ?;
// 		`, userID, password)
// 	if err != nil {
// 		return fmt.Errorf("Error while changing password! Please try again")
// 	}
// 	defer row.Close()
// 	if row.Next() {
// 		_, err := db.Exec(`
// 			UPDATE users
// 			SET password = ?
// 			WHERE user_id = ?;
// 			`, newPassword, userID)
// 		if err != nil {
// 			return fmt.Errorf("Error while changing password! Please try again")
// 		}
// 	} else {
// 		return fmt.Errorf("Wrong password! Plase try again")
// 	}

// 	return nil
// }
