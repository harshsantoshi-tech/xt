package handlers

import (
	"encoding/json"
	"errors"
	"expense-tracker/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// GetRoot is a simple handler for the root route.
func GetRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the complete API!")
}

func GetChatListPayload() models.ChatResponse {

	chats := []models.ChatEntry{
		{
			Contact: models.Contact{
				ID:        12,
				FirstName: "Vikrant",
				FullName:  "Vikrant",
				Phone:     "+917037306853",
			},
			Status:   "OPEN",
			Platform: "WHATSAPP",
			Assignment: models.Assignment{
				AssignedToType: "USER",
				AssignedToID:   intPtr(6),
				AssignedToName: "test 2",
			},
			LastMessage: models.LastMessage{
				ID:          340,
				Preview:     "Testingg",
				PreviewType: "text",
				Timestamp:   "2025-12-27T13:59:58+00:00",
				Direction:   "INCOMING",
				IsRead:      true,
			},
			UnreadCount: 0,
		},
		{
			Contact: models.Contact{
				ID:        15,
				FirstName: "Harsh",
				FullName:  "Santoshi",
				Phone:     "+917037306853",
			},
			Status:   "CLOSED",
			Platform: "WHATSAPP",
			Assignment: models.Assignment{
				AssignedToType: "USER",
				AssignedToID:   intPtr(6),
				AssignedToName: "test 2",
			},
			LastMessage: models.LastMessage{
				ID:          340,
				Preview:     "Production",
				PreviewType: "text",
				Timestamp:   "2025-12-27T13:59:58+00:00",
				Direction:   "INCOMING",
				IsRead:      true,
			},
			UnreadCount: 0,
		},
	}
	//chats = services.GetChatList(userID)

	// 2. Wrap in the top-level envelope
	return models.ChatResponse{
		Type:  "chat_list",
		Chats: chats,
	}
}

func intPtr(i int) *int { return &i }

func SendMessage(userID int64, payload json.RawMessage) (*models.MessageResponse, error) {

	var req models.SendMessagePayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, errors.New("invalid send_message payload")
	}

	if strings.TrimSpace(req.Message) == "" {
		return nil, errors.New("message cannot be empty")
	}

	//msg, err := services.SaveMessage(userID, req.ChatID, req.Message)
	//if err != nil {
	//	return nil, err
	//}

	// Broadcast to chat participants
	//services.BroadcastToChat(req.ChatID, models.MessageResponse{
	//	Type:    "new_message",
	//	Payload: msg,
	//})

	return &models.MessageResponse{
		Type:    "message_sent",
		Payload: "msg",
	}, nil
}

func BroadcastTyping(userID int64, payload json.RawMessage) {

	var req models.TypingPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return
	}

	//services.BroadcastToChat(req.ChatID, models.MessageResponse{
	//	Type: "typing",
	//	Payload: map[string]interface{}{
	//		"user_id": userID,
	//	},
	//})
}

func MarkMessageRead(userID int64, payload json.RawMessage) {

	var req models.ReadReceiptPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		return
	}

	//services.MarkMessageRead(userID, req.MessageID)

	//services.BroadcastMessageRead(req.MessageID, userID)
}
