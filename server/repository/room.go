package repository

import (
	"chat/domain"
	"database/sql"
	"log"

	"github.com/pkg/errors"
)

type roomRepository struct {
	DB *sql.DB
}

type RoomRepository interface {
	FindAllRooms() ([]*domain.Room, error)
	FindRoomByID(id int) (*domain.Room, error)
	CreateRoom(name string) (*domain.Room, error)
	// UpdateRoom()
	DeleteRoom(id int) error
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{
		DB: db,
	}
}

func (rr *roomRepository) FindAllRooms() ([]*domain.Room, error) {
	rows, err := rr.DB.Query("SELECT id, name FROM rooms")
	if err != nil {
		err = errors.Wrap(err, "failed to exec query")
		return nil, err
	}

	var rooms []*domain.Room
	for rows.Next() {
		var room domain.Room
		if err = rows.Scan(&room.ID, &room.Name); err != nil {
			log.Printf("failed to scan data: %v", err)
			continue
		}
		rooms = append(rooms, &room)
	}
	return rooms, nil
}

func (rr *roomRepository) FindRoomByID(id int) (*domain.Room, error) {
	row := rr.DB.QueryRow("SELECT id, name FROM rooms WHERE id=?", id)
	var room domain.Room
	if err := row.Scan(&room.ID, &room.Name); err != nil {
		if err == sql.ErrNoRows {
			return &room, nil
		}
		err = errors.Wrap(err, "failed to scan data")
		return nil, err
	}
	return &room, nil
}

func (rr *roomRepository) CreateRoom(name string) (*domain.Room, error) {
	result, err := rr.DB.Exec("INSERT INTO rooms(name) VALUES (?)", name)
	if err != nil {
		err = errors.Wrap(err, "failed to insert db")
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("failed to get new id: %v", err)
		err = errors.Wrap(err, "failed to get new id")
		return nil, err
	}
	return &domain.Room{
		ID:   int(id),
		Name: name,
	}, nil
}

func (rr *roomRepository) DeleteRoom(id int) error {
	_, err := rr.DB.Exec("DELETE FROM rooms WHERE id=?", id)
	if err != nil {
		err = errors.Wrap(err, "failed to delete")
		return err
	}
	return nil
}
