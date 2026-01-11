package controller

import (
	"expense-tracker/configs"
	"expense-tracker/models"
	"expense-tracker/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateRoomController(c echo.Context) error {
	var req models.CreateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status: "BAD_REQUEST",
			Message: "Invalid request",
			Code: http.StatusBadRequest,
		})
	}

	userId := c.Get("user_id").(int64)

	memberMap := make(map[int64]bool)

	memberMap[userId] = true

	for _, id := range req.Members {
		memberMap[id] = true
	}

	var uniqueMembers []int64
	for id := range memberMap {
		uniqueMembers = append(uniqueMembers, id)
	}

	if len(uniqueMembers) == 0  || len(uniqueMembers) == 1 {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status: "BAD_REQUEST",
			Message: "Invalid request as unique members are less than 2",
			Code: http.StatusBadRequest,
		})
	}

	if !req.IsGroup{
		req.Name = ""
	}

	roomID, existed, err := services.HandleCreateRoom(
		configs.AppConfig.DB,
		req.Name,
		req.IsGroup,
		uniqueMembers,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseModel{
			Status:  "INTERNAL_SERVER_ERROR",
			Message: fmt.Sprintf("error : %v", err),
			Code:    http.StatusInternalServerError,
		})
	}

	status := http.StatusCreated
	msg := "Room created successfully"
	if existed {
		status = http.StatusOK
		msg = "Room already exists"
	}

	return c.JSON(status, map[string]interface{}{
		"message": msg,
		"room_id": roomID,
	})
}