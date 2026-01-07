package controller

import (
	"encoding/json"
	"expense-tracker/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gorilla/websocket"
)

type SocketEvent struct {
	Type string `json:"type"`
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow Flutter connections
	},
}

func WsController(c echo.Context) error {

	userId, ok := c.Get("user_id").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return err
	}
	defer conn.Close()

	fmt.Println("Flutter client connected via Echo! ", userId)

	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		var event SocketEvent
		if err := json.Unmarshal(message, &event); err != nil {
			_ = conn.WriteJSON(map[string]string{
				"error": "invalid payload",
			})
			continue
		}

		switch event.Type {

		case "chat_list":
			res := handlers.GetChatListPayload()
			if err := conn.WriteJSON(res); err != nil {
				log.Println("Write error:", err)
				break
			}

		//case "send_message":
		//	res, err := handlers.SendMessage(event.Payload)
		//	if err != nil {
		//		_ = conn.WriteJSON(map[string]string{
		//			"type":  "error",
		//			"error": err.Error(),
		//		})
		//		continue
		//	}
		//	_ = conn.WriteJSON(res)
		//
		//case "typing":
		//	handlers.BroadcastTyping(event.Payload)
		//
		//case "read_receipt":
		//	handlers.MarkMessageRead(event.Payload)

		case "ping":
			_ = conn.WriteMessage(msgType, []byte(`{"type":"pong"}`))

		default:
			_ = conn.WriteJSON(map[string]string{
				"type":  "error",
				"error": "unknown event type",
			})
		}
	}

	log.Printf("webSocket disconnected user_id %d ", userId)
	return nil
}
