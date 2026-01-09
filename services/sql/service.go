package sql

import (
	"database/sql"
	"errors"
	"expense-tracker/configs"
	"expense-tracker/models"
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
			"UPDATE users SET password_hash = ?, created_at = ? WHERE email = ? AND status = 'pending'",
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

	query := "UPDATE users SET status = ? WHERE email = ? AND status = 'pending'"
	result, err := configs.AppConfig.DB.Exec(query,status, email)

	if err != nil {
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
	if status == "offline" {
		query = `UPDATE users SET status = ?, last_seen_at = CURRENT_TIMESTAMP WHERE id = ?`
	} else {
		query = `UPDATE users SET status = ? WHERE id = ?`
	}

	_, err := configs.AppConfig.DB.Exec(query, status, userID)
	return err
}

func GetChatListPayload(currentUserID int64, limit int, offset int) ([]models.ChatListItem, error) {

	rows, err := configs.AppConfig.DB.Query(GET_CHAT_LIST_PAGINATED, currentUserID, currentUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatList []models.ChatListItem
	for rows.Next() {
		var item models.ChatListItem
		var otherID sql.NullInt64
		var otherName, roomName, lastMsg, lastMsgType, lastSenderName sql.NullString
		var otherLastSeenAt time.Time

		err := rows.Scan(
			&item.RoomID, &roomName, &item.IsGroup, &item.LastMessageAt,
			&lastMsg, &lastMsgType, &item.LastSenderID, &lastSenderName,
			&item.UnreadCount, &otherID, &otherName,&otherLastSeenAt,
		)
		if err != nil {
			continue
		}

		if roomName.Valid { item.RoomName = &roomName.String }
		item.LastMessage = lastMsg.String
		item.MessageType = lastMsgType.String
		item.LastSenderName = lastSenderName.String

		if !item.IsGroup && otherID.Valid {
			item.OtherUserID = otherID.Int64
			item.OtherUserName = otherName.String
			item.OtherLastSeenAt = otherLastSeenAt
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

	updateRoom := `UPDATE rooms SET last_message_at = CURRENT_TIMESTAMP WHERE id = ?`
	if _, err := tx.Exec(updateRoom, roomID); err != nil {
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