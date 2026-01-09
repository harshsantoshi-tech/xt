package models

import (
	"encoding/json"
	"time"
)

type SocketEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type SocketResponse struct {
	Type string `json:"type"`
	Message interface{} `json:"message"`
}
type ChatMessage struct {
	ID          int64     `json:"id"`
	RoomID      int64     `json:"room_id"`
	SenderID    int64     `json:"sender_id"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type ChatListItem struct {
	RoomID        int64     `json:"room_id"`
	RoomName      *string   `json:"room_name"` // Pointer because it can be null for 1-on-1
	IsGroup       bool      `json:"is_group"`
	LastMessage   string    `json:"last_message"`
	MessageType   string    `json:"message_type"`
	LastSenderID  int64     `json:"last_sender_id"`
	LastSenderName string   `json:"last_sender_name"`
	LastMessageAt time.Time `json:"last_message_at"`
	UnreadCount   int       `json:"unread_count"`

	// Dynamic field for 1-on-1 chats
	OtherUserName string    `json:"other_user_name,omitempty"`
	OtherUserID   int64     `json:"other_user_id,omitempty"`
	OtherUserAvatar string `json:"other_user_avatar,omitempty"`
	OtherLastSeenAt time.Time `json:"other_last_seen_at,omitempty"`
}
// ChatResponse is the top-level envelope
type ChatResponse struct {
	Chats     []ChatListItem `json:"chats"`
	Limit     int         `json:"limit"`
	Offset    int         `json:"offset"`
	Timestamp string      `json:"timestamp"`
}

type FetchMessagesResponse struct {
	RoomID   int64         `json:"room_id"`
	Messages []ChatMessage `json:"messages"`
	Offset   int           `json:"offset"`
	Limit    int           `json:"limit"`
}

type ChatEntry struct {
	Contact     Contact     `json:"contact"`
	Status      string      `json:"status"`
	Platform    string      `json:"platform"`
	Assignment  Assignment  `json:"assignment"`
	LastMessage LastMessage `json:"last_message"`
	UnreadCount int         `json:"unread_count"`
}

type Contact struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
}

type Assignment struct {
	AssignedToType string `json:"assigned_to_type"`
	AssignedToID   *int   `json:"assigned_to_id"` // Pointer handles 'null'
	AssignedToName string `json:"assigned_to_name"`
}

type LastMessage struct {
	ID          int    `json:"id"`
	Preview     string `json:"preview"`
	PreviewType string `json:"preview_type"`
	Timestamp   string `json:"timestamp"`
	Direction   string `json:"direction"`
	AuthorType  string `json:"author_type"`
	Platform    string `json:"platform"`
	CreatedBy   string `json:"created_by"`
	ExpiresAt   int64  `json:"expires_at"`
	IsRead      bool   `json:"is_read"`
}

type ChatListResponse struct {
	Type  string      `json:"type"`
	Chats interface{} `json:"chats"`
}

type MessageResponse struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type SendMessagePayload struct {
	ChatID  int64  `json:"chat_id"`
	Message string `json:"message"`
}

type TypingPayload struct {
	ChatID int64 `json:"chat_id"`
}

type ReadReceiptPayload struct {
	MessageID int64 `json:"message_id"`
}
