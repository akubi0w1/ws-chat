package repository

import (
	"chat/domain"
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"
)

type messageRepository struct {
	DB *sql.DB
}

type MessageRepository interface {
	FindMessagesByRoomID(roomID int) ([]*domain.Message, error)
	FindMessageByID(id int) (*domain.Message, error)
	CreateMessage(body string, roomID, userID int, createdAt int64) (int, error)
	DeleteMessageByID(id int) error
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{
		DB: db,
	}
}

func (mr *messageRepository) FindMessagesByRoomID(roomID int) ([]*domain.Message, error) {
	rows, err := mr.DB.Query(`
		SELECT m.id, m.body, m.created_at, u.id, u.user_id, u.name, r.id, r.name
		FROM messages as m
		INNER JOIN users as u
		ON u.id = m.user_id
		INNER JOIN rooms as r
		ON r.id = m.room_id
		WHERE m.room_id=?
	`, roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*domain.Message{}, nil
		}
		return nil, err
	}
	var msgs []*domain.Message
	for rows.Next() {
		var msg domain.Message
		var user domain.User
		var room domain.Room
		var createdAt string
		if err = rows.Scan(&msg.ID, &msg.Body, &createdAt, &user.ID, &user.UserID, &user.Name, &room.ID, &room.Name); err != nil {
			log.Printf("failed to scan: %v", err)
			continue
		}
		t, _ := time.Parse("2006-01-02 15:04:05", createdAt)
		msg.Sender = &user
		msg.Room = &room
		msg.CreatedAt = t.Unix()
		msgs = append(msgs, &msg)
	}
	return msgs, nil
}

func (mr *messageRepository) FindMessageByID(id int) (*domain.Message, error) {
	row := mr.DB.QueryRow(`
		SELECT m.id, m.body, m.created_at, u.id, u.user_id, u.name, r.id, r.name
		FROM messages as m
		INNER JOIN users as u
		ON u.id = m.user_id
		INNER JOIN rooms as r
		ON r.id = m.room_id
		WHERE m.id=?
	`, id)
	var msg domain.Message
	var user domain.User
	var room domain.Room
	var createdAt string
	if err := row.Scan(&msg.ID, &msg.Body, &createdAt, &user.ID, &user.UserID, &user.Name, &room.ID, &room.Name); err != nil {
		if err == sql.ErrNoRows {
			return &msg, nil
		}
		err = errors.Wrap(err, "failed to get message")
		return nil, err
	}
	t, _ := time.Parse("2006-01-02 15:04:05", createdAt)
	msg.CreatedAt = t.Unix()
	msg.Sender = &user
	msg.Room = &room
	return &msg, nil
}

func (mr *messageRepository) CreateMessage(body string, roomID, userID int, createdAt int64) (int, error) {
	result, err := mr.DB.Exec(`
		INSERT INTO messages(body, user_id, room_id, created_at) VALUES (?,?,?,?)
	`, body, userID, roomID, time.Unix(createdAt, 0))
	if err != nil {
		err = errors.Wrap(err, "failed to insert db")
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		err = errors.Wrap(err, "failed to get new id")
		return 0, err
	}
	return int(id), nil
}

func (mr *messageRepository) DeleteMessageByID(id int) error {
	_, err := mr.DB.Exec("DELETE FROM messages WHERE id=?", id)
	if err != nil {
		err = errors.Wrap(err, "failed to delete data")
		return err
	}
	return nil
}
