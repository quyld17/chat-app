package messages

import (
	"database/sql"
	"time"
)

type Messages struct {
	Id        int       `json:"id"`
	RoomId    int       `json:"room_id"`
	UserId    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func StoreMessage(db *sql.DB, roomId, senderId int, message string) error {
	_, err := db.Exec(
		`INSERT INTO messages (room_id, user_id, content, created_at) 
		VALUES (?, ?, ?, ?);`,
		roomId, senderId, message, time.Now())
	return err
}

func GetChatHistory(db *sql.DB, roomId int) ([]Messages, error) {
	rows, err := db.Query(`
		SELECT id, room_id, sender_id, message, created_at 
		FROM messages 
		WHERE room_id = ? 
		ORDER BY created_at ASC`,
		roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Messages
	for rows.Next() {
		var msg Messages
		err := rows.Scan(&msg.Id, &msg.RoomId, &msg.UserId, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
