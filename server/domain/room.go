package domain

import (
	"database/sql"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Room struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	// websocket用のあれ
	Clients    map[*Client]bool
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Room *Room
	User *User
	Conn *websocket.Conn
	Send chan *Message
}

// TODO: 再定義はきしょい
type messageResponse struct {
	ID     int           `json:"id"`
	Body   string        `json:"body"`
	Room   *roomResponse `json:"room"`
	Sender *userResponse `json:"sender"`
}

// TODO: 再定義はきしょい
type roomResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TODO: 再定義はきしょい
type userResponse struct {
	ID     int    `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func NewRoom(id int, name string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// TODO: これ以下、serviceみたいに、層作った方がいい。
// repository呼びたい時にimport cycleしてしまう
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}
		case message := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.Clients, client)
				}
			}
		}
	}
}

// client -> server
func (c *Client) ReadMessage() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var message MessageRequest
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			log.Printf("failed to read message: %v", err)
			break
		}
		msg := &Message{
			Body:      message.Body,
			CreatedAt: time.Now().Unix(),
			Room:      c.Room,
			Sender:    c.User,
		}
		c.Room.Broadcast <- msg
	}
}

// server -> client
// func (c *Client) WriteMessage(repo repository.MessageRepository) {
func (c *Client) WriteMessage() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	for {
		message := <-c.Send
		// TODO: db登録のタイミングが悪すぎる...?
		id, err := insertMessage(message)
		if err != nil {
			log.Printf("failed to insert db: %v", err)
			continue
		}
		msg := &messageResponse{
			ID:     int(id),
			Body:   message.Body,
			Room:   &roomResponse{ID: message.Room.ID, Name: message.Room.Name},
			Sender: &userResponse{ID: message.Sender.ID, UserID: message.Sender.UserID, Name: message.Sender.Name},
		}

		err = c.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("failed to write message: %v", err)
			break
		}
	}
}

// TODO: なんかいい方法ないかなこれ。
// repositoryをimportするとcycleしてしまうんよな。
// どっかにコネクション持たせておくとかしないとなあ
func insertMessage(msg *Message) (int, error) {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3307)/chat")
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	result, err := db.Exec(`
		INSERT INTO messages(body, user_id, room_id, created_at) VALUES (?,?,?,?)
	`, msg.Body, msg.Sender.ID, msg.Room.ID, time.Unix(msg.CreatedAt, 0))
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
