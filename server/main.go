package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"

	"chat/aggregation"
	"chat/handler"
)

// var upgrader = websocket.Upgrader{}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("[info] finish setting log")
}

func main() {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3307)/chat")
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	if err := aggregation.InitRoomAggregation(db); err != nil {
		log.Fatalf("failed init: %v", err)
	}

	roomHandler := handler.NewRoomHandler(db)
	userHandler := handler.NewUserHandler(db)
	msgHandler := handler.NewMessageHandler(db)

	// routing
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/ping", HealthHandler)

	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	r.HandleFunc("/rooms", roomHandler.GetAllRooms).Methods("GET")
	r.HandleFunc("/rooms", roomHandler.CreateRoom).Methods("POST")

	r.HandleFunc("/rooms/{roomID:[0-9]+}/messages", msgHandler.GetMessagesByRoomID).Methods("GET")
	r.HandleFunc("/rooms/{roomID:[0-9]+}/messages", msgHandler.CreateMessage).Methods("POST")

	// websocketコネクションの接続用
	r.HandleFunc("/ws/{id:[0-9]+}", handler.StartConnection)

	http.ListenAndServe(":8080", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
