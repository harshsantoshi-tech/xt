package sql

import (
	"database/sql"
	"errors"
	"expense-tracker/configs"
	"expense-tracker/models"
	"fmt"
	"github.com/labstack/gommon/log"
	"time"
)

func UpsertUser(email , hashedPassword string)error{
	var status string
	err := configs.AppConfig.DB.QueryRow("SELECT status FROM users WHERE email = ?", email).Scan(&status)

	if err == nil {
		// User exists - check their status
		if status != "pending" {
			return errors.New("user already registered with this email")
		}
		// If they are pending, update their password (in case they changed it on retry)
		_, err = configs.AppConfig.DB.Exec(
			"UPDATE users SET password_hash = ?, updated_at = ? WHERE email = ? AND status = 'pending'",
			string(hashedPassword), time.Now(), email,
		)
	} else if errors.Is(err, sql.ErrNoRows) {
		// New user - perform initial insert
		_, err = configs.AppConfig.DB.Exec(
			"INSERT INTO users (email, password_hash, status, created_at) VALUES (?, ?, 'pending', ?)",
			email, string(hashedPassword), time.Now(),
		)
	}

	return err
}


func UpdateUserStatus(email string, status string)error{

	query := "UPDATE users SET status = ? , last_seen_at = ? WHERE email = ? AND status = 'pending'"
	result, err := configs.AppConfig.DB.Exec(query,status,time.Now(), email)

	if err != nil {
		log.Error("UpdateUserStatus.Exec err: %v", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		err = errors.New("registration failed: user not found or already verified")
	}
	return err
}

func UpdateUserOnlineStatus(userID int64, status string) error {
	var query string
	query = `UPDATE users SET status = ?, last_seen_at = ? WHERE id = ?`

	_, err := configs.AppConfig.DB.Exec(query, status,time.Now(), userID)
	return err
}

func GetChatListPayload(currentUserID int64, limit int, offset int) ([]models.ChatListItem, error) {

	rows, err := configs.AppConfig.DB.Query(GET_CHAT_LIST_PAGINATED, currentUserID, currentUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var chatList []models.ChatListItem
	for rows.Next() {
		var item models.ChatListItem

		var lastMsgAt, otherSeenAt time.Time

		err := rows.Scan(
			&item.RoomID,
			&item.RoomName,
			&item.IsGroup,
			&lastMsgAt,
			&item.LastMessage,
			&item.MessageType,
			&item.LastSenderID,
			&item.LastSenderName,
			&item.UnreadCount,
			&item.OtherUserID,
			&item.OtherUserName,
			&otherSeenAt,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		item.LastMessageAt = lastMsgAt

		if !item.IsGroup && !otherSeenAt.IsZero() {
			item.OtherLastSeenAt = &otherSeenAt
		}

		chatList = append(chatList, item)
	}

	return chatList, nil
}

func SaveMessage(roomID int64, senderID int64, content string, msgType string) (models.ChatMessage, error) {
	tx, err := configs.AppConfig.DB.Begin()
	if err != nil {
		return models.ChatMessage{}, err
	}

	query := `INSERT INTO messages (room_id, sender_id, content, message_type) VALUES (?, ?, ?, ?)`
	res, err := tx.Exec(query, roomID, senderID, content, msgType)
	if err != nil {
		tx.Rollback()
		return models.ChatMessage{}, err
	}

	id, _ := res.LastInsertId()

	updateRoom := `UPDATE rooms SET last_message_at = ? WHERE id = ?`
	if _, err := tx.Exec(updateRoom,time.Now(), roomID); err != nil {
		tx.Rollback()
		return models.ChatMessage{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.ChatMessage{}, err
	}

	return models.ChatMessage{
		ID: id,
		RoomID: roomID,
		SenderID: senderID,
		Content: content,
		MessageType: msgType,
		CreatedAt: time.Now(),
	}, nil
}

func GetRoomMembers(roomId int64) ([]int64, error) {
	var users []int64
	rows, err := configs.AppConfig.DB.Query("SELECT user_id FROM room_members WHERE room_id = ?", roomId)
	if err != nil {
		return nil , err
	}
	defer rows.Close()

	for rows.Next() {
		var memberID int64
		errS := rows.Scan(&memberID)
		if errS != nil{
			continue
		}
		users = append(users, memberID)
	}
	return users ,nil
}

func MarkMessagesAsRead(roomID int64, userID int64, lastMessageID int64) error {
	query := `UPDATE room_members SET last_read_message_id = ? 
              WHERE room_id = ? AND user_id = ?`
	_, err := configs.AppConfig.DB.Exec(query, lastMessageID, roomID, userID)
	return err
}

func GetMessagesByRoom(roomID int64, limit int, offset int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage

	rows, err := configs.AppConfig.DB.Query(GET_MESSAGES_BY_ROOM, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.ChatMessage
		err := rows.Scan(
			&m.ID,
			&m.RoomID,
			&m.SenderID,
			&m.Content,
			&m.MessageType,
			&m.CreatedAt,
		)
		if err == nil {
			messages = append(messages, m)
		}
	}

	return messages, nil
}

func SaveUserSession(userID int64, refreshToken string, expiresAt time.Time) error {
	query := `
		INSERT INTO user_sessions (user_id, refresh_token, expires_at) 
		VALUES (?, ?, ?)
	`

	_, err := configs.AppConfig.DB.Exec(query, userID, refreshToken, expiresAt)
	if err != nil {
		return err
	}

	return nil
}