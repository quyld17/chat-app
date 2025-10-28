package status

import (
	"database/sql"
	"fmt"
)

type Status struct {
	UserId int `json:"user_id"`
}

func Update(db *sql.DB, userID int) error {
	query := `	INSERT INTO status (user_id) 
				VALUES (?) 
				ON DUPLICATE KEY UPDATE user_id = user_id;`

	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("Failed to update status: %v", err)
	}

	return nil
}
func Remove(db *sql.DB, userID int) error {
	query := `	DELETE FROM status 
				WHERE user_id = ?;`

	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("Failed to update status")
	}

	return nil
}
