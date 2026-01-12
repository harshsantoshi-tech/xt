package handlers

import (
	"encoding/json"
	"expense-tracker/constants"
	"expense-tracker/models"
	"expense-tracker/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

// GetRoot is a simple handler for the root route.
func GetRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the complete API!")
}

// WritePump handles sending messages FROM the server TO the Flutter client
func WritePump(client *services.Client) {
	defer client.Conn.Close()
	for message := range client.Send {
		if err := client.Conn.WriteJSON(message); err != nil {
			return
		}
	}
}

// ReadPump handles receiving messages FROM the Flutter client TO the server
func ReadPump(client *services.Client) {
	defer func() {
		services.BroadcastUserStatus(client.UserID,constants.OFFLINE)
		services.Hub.Unregister(client.UserID)
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("User %d disconnected", client.UserID)
			break
		}

		var event models.SocketEvent
		if err = json.Unmarshal(message, &event); err != nil {
			log.Error("Error unmarshalling message ",message)
			continue
		}

		handleSocketEvent(client, event)
	}
}

func handleSocketEvent(client *services.Client, event models.SocketEvent) {
	switch event.Type {

	case "chat_list":

		res , err := services.GetChatList(client.UserID ,event.Payload)
		if err != nil {
			client.Send <- models.SocketResponse{Type: "error", Message: err.Error()}
		} else {
			client.Send <- models.SocketResponse{Type: "chat_list", Message: res}
		}

	case "send_message":
		msg , err := services.HandleSendMessage(client, event.Payload)
		if err != nil {
			client.Send <- models.SocketResponse{Type: "error", Message: err.Error()}
		} else {
			client.Send <- models.SocketResponse{Type: "send_message", Message: msg}
		}

	case "typing":
		err := services.HandleTypingEvent(client , event.Payload)
		if err != nil {
			client.Send <- models.SocketResponse{Type: "error", Message: err.Error()}
		} else {
			client.Send <- models.SocketResponse{Type: "typing", Message: nil}
		}

	case "read_receipt":
		err := services.HandleReadReceipt(client, event.Payload)
		if err != nil {
			client.Send <- models.SocketResponse{Type: "error", Message: err.Error()}
		} else {
			client.Send <- models.SocketResponse{Type: "read_receipt", Message: nil}
		}

	case "fetch_messages":
		res, err := services.HandleFetchMessages(client, event.Payload)
		if err != nil {
			client.Send <- models.SocketResponse{Type: "error", Message: err.Error()}
		} else {
			client.Send <- models.SocketResponse{Type: "fetch_messages", Message: res}
		}

	case "ping":
		client.Send <- map[string]string{"type": "pong"}
	}
}