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
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func Save(db *sql.DB, roomId, senderId int, message string) error {
	_, err := db.Exec(`
		INSERT INTO messages (room_id, user_id, content) 
		VALUES (?, ?, ?);`,
		roomId, senderId, message)
	return err
}

func GetChatHistory(db *sql.DB, roomId, limit int) ([]Messages, error) {
	rows, err := db.Query(`
		SELECT 
			messages.id, 
			messages.room_id, 
			messages.user_id, 
			users.username, 
			messages.content, 
			messages.created_at 
		FROM messages
		INNER JOIN users ON messages.user_id = users.id
		WHERE messages.room_id = ? 
		ORDER BY messages.created_at DESC
		LIMIT ?;`,
		roomId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Messages
	for rows.Next() {
		var msg Messages
		err := rows.Scan(&msg.Id, &msg.RoomId, &msg.UserId, &msg.Username, &msg.Content, &msg.CreatedAt)
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
