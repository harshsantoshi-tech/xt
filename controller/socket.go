package controller

import (
	"expense-tracker/handlers"
	"expense-tracker/services"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WsController(c echo.Context) error {
	userId := c.Get("user_id").(int64)

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := &services.Client{
		UserID: userId,
		Conn:   conn,
		Send:   make(chan interface{}, 256),
	}

	services.Hub.Register(client)
	go services.BroadcastUserStatus(client.UserID,"online")
	// Start concurrent workers
	go handlers.WritePump(client)
	handlers.ReadPump(client)

	return nil
}