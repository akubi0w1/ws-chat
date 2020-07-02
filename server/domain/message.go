package domain

import (
	"time"
)

type Message struct {
	ID        int
	Body      string
	CreatedAt int64
	Room      *Room
	Sender    *User
}

type MessageRequest struct {
	Body     string `json:"body"`
	SenderID int    `json:"sender_id"`
}

func NewMessage(body string, sender *User, room *Room) *Message {
	return &Message{
		Body:      body,
		CreatedAt: time.Now().Unix(),
		Sender:    sender,
		Room:      room,
	}
}
