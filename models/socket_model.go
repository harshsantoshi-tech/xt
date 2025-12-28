package models



// ChatResponse is the top-level envelope
type ChatResponse struct {
	Type      string      `json:"type"`
	Chats     []ChatEntry `json:"chats"`
	Limit     int         `json:"limit"`
	Offset    int         `json:"offset"`
	Timestamp string      `json:"timestamp"`
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