package services

import (
	"database/sql"
	"fmt"
)

func HandleCreateRoom(db *sql.DB, name string, isGroup bool, members []int64) (int64, bool, error) {

	if !isGroup && len(members) == 2 {
		var existingRoomID int64
		query := `
            SELECT r.id FROM rooms r
            JOIN room_members rm ON r.id = rm.room_id
            WHERE r.is_group = 0 AND rm.user_id IN (?, ?)
            GROUP BY r.id HAVING COUNT(DISTINCT rm.user_id) = 2`

		err := db.QueryRow(query, members[0], members[1]).Scan(&existingRoomID)
		if err == nil {
			return existingRoomID, true, nil // true = already existed
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, false, err
	}

	res, err := tx.Exec("INSERT INTO rooms (name, is_group) VALUES (?, ?)", name, isGroup)
	if err != nil {
		tx.Rollback()
		return 0, false, err
	}

	roomID, _ := res.LastInsertId()

	for _, userID := range members {
		_, err := tx.Exec("INSERT INTO room_members (room_id, user_id) VALUES (?, ?)", roomID, userID)
		if err != nil {
			tx.Rollback()
			return 0, false, fmt.Errorf("failed to add member %d: %w", userID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, false, err
	}

	return roomID, false, nil
}
