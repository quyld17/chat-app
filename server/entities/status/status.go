package status

import (
	"database/sql"
	"fmt"
)

type Status struct {
	UserId   int  `json:"user_id"`
	IsOnline bool `json:"is_online"`
}

func Update(db *sql.DB, userID int) error {
	query := `	INSERT INTO status (user_id) 
				VALUES (?);`

	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to update status for user %s: %w", userID, err)
	}

	return nil
}

func Remove(db *sql.DB, userID int) error {
	query := `	DELETE FROM status 
				WHERE user_id = ?;`

	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to remove status for user %s: %w", userID, err)
	}

	return nil
}
