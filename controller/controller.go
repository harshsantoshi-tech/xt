package controller

import (
	"encoding/json"
	"expense-tracker/handlers"
	"expense-tracker/services"
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/gorilla/websocket"
)

type SocketEvent struct {
	Type string          `json:"type"`
}

func HomePageController(c echo.Context) error {

	res, err := handlers.HomePageHandler(1)
	if err != nil {
		log.Println("Error giving home page data ", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}


// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow Flutter connections
	},
}

func WsController(c echo.Context) error {
	// 1. Upgrade the Echo connection to a WebSocket connection
	// We pass the raw ResponseWriter and Request from Echo
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return err // Echo will handle this error
	}
	defer conn.Close()

	token := c.QueryParam("token")
	secret := os.Getenv("JWT_SECRET")
	userId, err := services.GetUserIDFromToken(token, secret)

	log.Println("token ",token ," userid " ,userId)
	fmt.Println("Flutter client connected via Echo!")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil { break }
	
		var event SocketEvent

		if err := json.Unmarshal(message, &event); err != nil { continue }
	
		// Dispatch based on Type
		switch event.Type {
		case "chat_list":
			res := handlers.GetChatListPayload()
            
            if err := conn.WriteJSON(res); err != nil {
                log.Println("Write error:", err)
            }
		default:
			fmt.Println("Unknown event type:", event.Type)
		}
	}

	return nil // Connection closed
}




