package handler

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"chat/aggregation"
	"chat/domain"
	"chat/repository"
	"chat/response"
)

type RoomHandler interface {
	GetAllRooms(w http.ResponseWriter, r *http.Request)
	// GetRoomByID(w http.ResponseWriter, r *http.Request)
	CreateRoom(w http.ResponseWriter, r *http.Request)
}

type roomHandler struct {
	repository.RoomRepository
}

func NewRoomHandler(db *sql.DB) RoomHandler {
	return &roomHandler{
		repository.NewRoomRepository(db),
	}
}

func (rh *roomHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := rh.RoomRepository.FindAllRooms()
	if err != nil {
		log.Printf("failed to get all rooms: %v", err)
		response.InternalServerError(w, err)
		return
	}

	var res RoomsResponse
	for _, room := range rooms {
		res.Rooms = append(res.Rooms, convToRoomResponse(room))
	}
	response.Success(w, res)
}

func (rh *roomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		response.BadRequest(w, err)
		return
	}
	var req CreateRoomRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("failed to unmarshal json: %v", err)
		response.InternalServerError(w, err)
		return
	}

	re, err := rh.RoomRepository.CreateRoom(req.Name)
	if err != nil {
		log.Printf("failed to create room: %v", err)
		response.InternalServerError(w, err)
		return
	}

	room := domain.NewRoom(re.ID, re.Name)
	aggregation.RoomAggregation[room.ID] = room

	response.Success(w, convToRoomResponse(room))
}

func convToRoomResponse(room *domain.Room) *RoomResponse {
	return &RoomResponse{
		ID:   room.ID,
		Name: room.Name,
	}
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type RoomResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RoomsResponse struct {
	Rooms []*RoomResponse `json:"rooms"`
}
