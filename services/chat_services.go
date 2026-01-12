package services

import (
	"encoding/json"
	"expense-tracker/configs"
	"expense-tracker/models"
	"expense-tracker/services/sql"
	"github.com/labstack/gommon/log"
	"time"
)

func GetChatList(userId int64, payload json.RawMessage) (models.ChatResponse, error) {

	var pagination models.Pagination
	err := json.Unmarshal(payload, &pagination)
	if err != nil {
		return models.ChatResponse{}, err
	}

	if pagination.Limit <= 0 {
		pagination.Limit = 20
	}

	chats, err := sql.GetChatListPayload(userId, pagination.Limit, pagination.Offset)
	if err != nil {
		log.Error("GetChats err: ", err, " ", userId)
		return models.ChatResponse{}, err
	}

	return models.ChatResponse{
		Chats:  chats,
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	}, nil
}

func HandleSendMessage(client *Client, payload json.RawMessage) (models.ChatMessage, error) {
	var req models.SendMessagesPayload

	if err := json.Unmarshal(payload, &req); err != nil {
		log.Print("Error unmarshaling send_message req:", err)
		return models.ChatMessage{}, err
	}

	msg, err := sql.SaveMessage(req.RoomID, client.UserID, req.Content, req.Type)
	if err != nil {
		log.Print("Error saving message:", err)
		return models.ChatMessage{}, err
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Print("Error marshaling message for broadcast:", err)
		return models.ChatMessage{}, err
	}

	event := models.SocketEvent{
		Type:    "new_message",
		Payload: msgBytes,
	}

	roomMembers, err := sql.GetRoomMembers(req.RoomID)

	if err != nil {
		return models.ChatMessage{}, err
	}

	for _, memberId := range roomMembers {
		if memberId != client.UserID {
			Hub.SendToUser(memberId, event)
		}
	}
	return msg, nil
}

func HandleTypingEvent(client *Client, payload json.RawMessage) error {
	var req models.TypingPayload

	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	broadcastPayload, _ := json.Marshal(map[string]interface{}{
		"room_id":   req.RoomID,
		"user_id":   client.UserID,
		"is_typing": req.IsTyping,
	})

	broadcastEvent := models.SocketEvent{
		Type:    "typing",
		Payload: broadcastPayload,
	}

	roomMembers, err := sql.GetRoomMembers(req.RoomID)

	if err != nil {
		return err
	}

	for _, memberId := range roomMembers {
		if memberId != client.UserID {
			Hub.SendToUser(memberId, broadcastEvent)
		}
	}
	return nil
}

func BroadcastUserStatus(userID int64, status string) {

	go sql.UpdateUserOnlineStatus(userID, status)

	event := models.SocketEvent{
		Type:    "user_status",
		Payload: nil,
	}

	query := ` SELECT DISTINCT user_id FROM room_members 
        WHERE room_id IN (SELECT room_id FROM room_members WHERE user_id = ?)
        AND user_id != ?`

	rows, err := configs.AppConfig.DB.Query(query, userID, userID)
	if err != nil {
		return
	}
	defer rows.Close()

	statusData, _ := json.Marshal(map[string]interface{}{
		"user_id":      userID,
		"status":       status,
		"last_seen_at": time.Now().Format(time.RFC3339),
	})
	event.Payload = statusData

	for rows.Next() {
		var contactID int64
		rows.Scan(&contactID)
		Hub.SendToUser(contactID, event)
	}
}

func HandleReadReceipt(client *Client, payload json.RawMessage) error {
	var req models.ReadReceiptPayload

	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}

	_ = sql.MarkMessagesAsRead(req.RoomID, client.UserID, req.LastMessageID)

	broadcastPayload, _ := json.Marshal(map[string]interface{}{
		"room_id":              req.RoomID,
		"reader_id":            client.UserID,
		"last_read_message_id": req.LastMessageID,
	})

	receiptEvent := models.SocketEvent{
		Type:    "read_receipt",
		Payload: broadcastPayload,
	}

	roomMembers, err := sql.GetRoomMembers(req.RoomID)

	if err != nil {
		return err
	}

	for _, memberId := range roomMembers {
		if memberId != client.UserID {
			Hub.SendToUser(memberId, receiptEvent)
		}
	}
	return nil
}

func HandleFetchMessages(client *Client, payload json.RawMessage) (models.FetchMessagesResponse, error) {
	var req struct {
		RoomID int64 `json:"room_id"`
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
	}

	if err := json.Unmarshal(payload, &req); err != nil {
		return models.FetchMessagesResponse{}, err
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}

	messages, err := sql.GetMessagesByRoom(req.RoomID, req.Limit, req.Offset)
	if err != nil {
		return models.FetchMessagesResponse{}, err
	}

	response := models.FetchMessagesResponse{
		RoomID:   req.RoomID,
		Messages: messages,
		Offset:   req.Offset,
		Limit:    req.Limit,
	}

	return response, nil
}
