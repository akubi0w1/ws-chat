package handler

import (
	"chat/aggregation"
	"chat/domain"
	"chat/repository"
	"chat/response"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var upgrader = websocket.Upgrader{}

func StartConnection(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	roomID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("failed to parse int: %v", err)
		response.BadRequest(writer, "failed to get roomID")
		return
	}
	room, ok := aggregation.RoomAggregation[roomID]
	if !ok {
		log.Printf("failed to find a room")
		response.BadRequest(writer, "failed to find a room")
		return
	}
	// userの取得
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3307)/chat")
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	cookie, err := request.Cookie("session_user_id")
	if err != nil {
		log.Printf("failed to get cookie: %v", err)
		response.BadRequest(writer, err)
		return
	}
	userID := cookie.Value
	row := db.QueryRow("SELECT id, user_id, name FROM users WHERE id=?", userID)
	var user domain.User
	if err := row.Scan(&user.ID, &user.UserID, &user.Name); err != nil {
		err = errors.Wrap(err, "failed to get user")
		response.InternalServerError(writer, err)
		return
	}

	repo := repository.NewMessageRepository(db)

	go room.Run()

	serveWebsocket(room, &user, repo, writer, request)
}

func serveWebsocket(room *domain.Room, user *domain.User, repo repository.MessageRepository, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade get request to a websocket: %v", err)
		// TODO: しかるべき措置...
		return
	}
	client := &domain.Client{Room: room, User: user, Conn: conn, Send: make(chan *domain.Message)}
	client.Room.Register <- client

	log.Println("serv websocket")
	go client.WriteMessage()
	go client.ReadMessage()
}
