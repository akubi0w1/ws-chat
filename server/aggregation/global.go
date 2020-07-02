package aggregation

import (
	"chat/domain"
	"database/sql"
	"log"

	"github.com/pkg/errors"
)

var RoomAggregation = make(map[int]*domain.Room)

func InitRoomAggregation(db *sql.DB) error {
	// roomを全件取得
	rows, err := db.Query("SELECT id, name FROM rooms")
	if err != nil {
		err = errors.Wrap(err, "failed to exec query")
		return err
	}

	for rows.Next() {
		var room domain.Room
		if err = rows.Scan(&room.ID, &room.Name); err != nil {
			log.Printf("failed to scan data: %v", err)
			continue
		}
		RoomAggregation[room.ID] = domain.NewRoom(room.ID, room.Name)
	}
	return nil
}
