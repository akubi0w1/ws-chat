package handler

import (
	"chat/domain"
	"chat/repository"
	"chat/response"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type messageHandler struct {
	repository.MessageRepository
}

type MessageHandler interface {
	GetMessagesByRoomID(w http.ResponseWriter, r *http.Request)
	CreateMessage(w http.ResponseWriter, r *http.Request)
}

func NewMessageHandler(db *sql.DB) MessageHandler {
	return &messageHandler{
		MessageRepository: repository.NewMessageRepository(db),
	}
}

func (mh *messageHandler) GetMessagesByRoomID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["roomID"])
	if err != nil {
		log.Printf("failed to parse int: %v", err)
		response.BadRequest(w, "failed to get roomID")
		return
	}

	msgs, err := mh.MessageRepository.FindMessagesByRoomID(roomID)
	if err != nil {
		log.Printf("failed to get messages: %v", err)
		response.InternalServerError(w, err)
		return
	}

	var res MessagesResponse
	for _, msg := range msgs {
		res.Messages = append(res.Messages, convToMessageResponse(msg))
	}
	response.Success(w, res)
}

func (mh *messageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["roomID"])
	if err != nil {
		log.Printf("failed to parse int: %v", err)
		response.BadRequest(w, "failed to get roomID")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		response.BadRequest(w, err)
		return
	}
	var req CreateMessageRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("failed to unmarshal json: %v", err)
		response.InternalServerError(w, err)
		return
	}
	createdAt := time.Now().Unix()
	id, err := mh.MessageRepository.CreateMessage(req.Body, roomID, req.SenderID, createdAt)
	if err != nil {
		log.Printf("failed to create data: %v", err)
		response.InternalServerError(w, err)
		return
	}
	// スマートにやりたいなあ...
	msg, err := mh.MessageRepository.FindMessageByID(id)
	if err != nil {
		log.Printf("failed to get new data: %v", err)
		response.InternalServerError(w, err)
		return
	}

	response.Success(w, convToMessageResponse(msg))
}

func convToMessageResponse(msg *domain.Message) *MessageResponse {
	return &MessageResponse{
		ID:        msg.ID,
		Body:      msg.Body,
		CreatedAt: msg.CreatedAt,
		Sender:    convToUserResponse(msg.Sender),
		Room:      convToRoomResponse(msg.Room),
	}
}

type CreateMessageRequest struct {
	Body     string `json:"body"`
	SenderID int    `json:"sender_id"`
}

type MessageResponse struct {
	ID        int           `json:"id"`
	Body      string        `json:"body"`
	CreatedAt int64         `json:"created_at"`
	Sender    *UserResponse `json:"sender"`
	Room      *RoomResponse `json:"room"`
}

type MessagesResponse struct {
	Messages []*MessageResponse `json:"messages"`
}
