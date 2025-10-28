package rooms

import (
	"database/sql"
	"fmt"
	"time"
)

type Rooms struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	IsGroup   bool      `json:"is_group"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatParticipants struct {
	Id     int `json:"id"`
	RoomId int `json:"room_id"`
	UserId int `json:"user_id"`
}

func GetId(db *sql.DB, receiverId, senderId int) (int, error) {
	var roomId int
	queryErr := db.QueryRow(`
		SELECT room_id
		FROM chat_participants
		WHERE user_id IN (?, ?)
		GROUP BY room_id
		HAVING COUNT(DISTINCT user_id) = 2;`,
		receiverId, senderId).Scan(&roomId)

	if queryErr != nil {
		return 0, queryErr
	}

	return roomId, nil
}

func Create(db *sql.DB, receiverId, senderId int) (int, error) {
	var roomId int
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	var commitErr error

	defer func() {
		if commitErr != nil {
			tx.Rollback()
			fmt.Printf("transaction rollback due to: %v\n", commitErr)
		} else {
			if commitErr = tx.Commit(); commitErr != nil {
				tx.Rollback()
				fmt.Printf("transaction commit failed: %v\n", commitErr)
			}
		}
	}()

	result, err := tx.Exec(`
			INSERT INTO rooms (name)
			VALUES ("");`)
	if err != nil {
		commitErr = fmt.Errorf("failed to insert into rooms: %v", err)
		return 0, commitErr
	}

	newRoomId, err := result.LastInsertId()
	if err != nil {
		commitErr = fmt.Errorf("failed to retrieve last insert id: %v", err)
		return 0, commitErr
	}

	_, err = tx.Exec(`
			INSERT INTO chat_participants (room_id, user_id)
			VALUES (?, ?), (?, ?)`,
		newRoomId, senderId, newRoomId, receiverId)
	if err != nil {
		commitErr = fmt.Errorf("failed to insert chat participants: %v", err)
		return 0, commitErr
	}

	roomId = int(newRoomId)

	return roomId, nil
}
