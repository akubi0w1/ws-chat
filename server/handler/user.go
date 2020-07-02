package handler

import (
	"chat/domain"
	"chat/repository"
	"chat/response"
	"strconv"

	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type userHandler struct {
	repository.UserRepository
}

type UserHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(db *sql.DB) UserHandler {
	return &userHandler{
		repository.NewUserRepository(db),
	}
}

// TODO: パスワードとは？
func (uh *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		response.BadRequest(w, err)
		return
	}
	var req LoginRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("failed to unmarshal json: %v", err)
		response.InternalServerError(w, err)
		return
	}
	user, err := uh.UserRepository.FindUsreByUserID(req.UserID)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		response.InternalServerError(w, err)
		return
	}
	if user == nil {
		log.Printf("failed to find user. userID=%v", req.UserID)
		response.BadRequest(w, "failed to find user")
		return
	}

	// TODO: 直書きやめろ
	cookie := &http.Cookie{
		Name:  "session_user_id",
		Value: strconv.Itoa(user.ID),
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	response.Success(w, convToUserResponse(user))
}

func (uh *userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.UserRepository.FindAllUsers()
	if err != nil {
		log.Printf("failed to get all users: %v", err)
		response.InternalServerError(w, err)
		return
	}
	var res UsersResponse
	for _, user := range users {
		res.Users = append(res.Users, convToUserResponse(user))
	}
	response.Success(w, res)
}

func (uh *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		response.BadRequest(w, err)
		return
	}
	var req CreateUserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("failed to unmarshal json: %v", err)
		response.InternalServerError(w, err)
		return
	}
	if req.UserID == "" || req.Name == "" {
		err = errors.New("user_id and name is require")
		response.BadRequest(w, err.Error())
		return
	}

	id, err := uh.UserRepository.CreateUser(req.UserID, req.Name)
	user := domain.NewUser(id, req.UserID, req.Name)
	response.Success(w, convToUserResponse(user))
}

func convToUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:     user.ID,
		UserID: user.UserID,
		Name:   user.Name,
	}
}

type LoginRequest struct {
	UserID string `json:"user_id"`
}

type CreateUserRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type UsersResponse struct {
	Users []*UserResponse `json:"users"`
}

type UserResponse struct {
	ID     int    `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}
