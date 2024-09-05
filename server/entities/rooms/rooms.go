package rooms

import (
	"database/sql"
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

func GetRoom(db *sql.DB, receiverId, senderId int) (int, error) {
	var roomId int

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.QueryRow(`
		SELECT room_id
		FROM ChatParticipants
		WHERE user_id IN (?, ?)
		GROUP BY room_id
		HAVING COUNT(DISTINCT user_id) = 2;`,
		receiverId, senderId).Scan(&roomId)

	if err == sql.ErrNoRows {
		result, err := tx.Exec(`
			INSERT INTO Rooms (created_at)
			VALUES (?, ?, ?);`,
			time.Now())
		if err != nil {
			return 0, err
		}

		newRoomId, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		_, err = tx.Exec(`
			INSERT INTO ChatParticipants (room_id, user_id)
			VALUES (?, ?), (?, ?)`,
			newRoomId, senderId, newRoomId, receiverId)
		if err != nil {
			return 0, err
		}

		roomId = int(newRoomId)
	} else if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return roomId, nil
}


