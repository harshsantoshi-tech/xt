package handlers

import (
	"expense-tracker/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetRoot is a simple handler for the root route.
func GetRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the complete API!")
}

// TextChangeHandler handles the /text/change endpoint, demonstrating the use of the new models.
func TextChangeHandler(c echo.Context) error {
	// Default values
	color := ""
	fontSize := float64(0)

	// Check for query parameters and override defaults if present
	if c.QueryParam("color") != "" {
		color = c.QueryParam("color")
	}

	if c.QueryParam("fontSize") != "" {
		if size, err := strconv.ParseFloat(c.QueryParam("fontSize"), 64); err == nil {
			fontSize = size
		}
	}

	completedTag := models.OfferWidgetTag{
		//Title: models.UikitText{
		//	Value: "Completed",
		//},
		Style: models.UikitStyle{
			Color:    color,
			FontSize: fontSize,
		},
		ButtonStyle: models.UikitStyle{
			Color:    "",
			FontSize: 0.0,
			Value:    "",
		},
	}

	res := models.Response{
		Response: models.ResponseModel{
			Status:  "SUCCESS",
			Message: "",
			Code:    http.StatusOK,
		},
		Compete: completedTag,
	}

	// Returns the response as JSON
	return c.JSON(http.StatusOK, res)
}

func HomePageHandler(userId int64) (models.HomePage, error) {

	//expenseItems := GetAllExpenseItemsForUser(userId)

	res := models.HomePage{
		CommonAppBar: models.UikitText{
			Title: "Xt",
			TitleStyle: models.UikitStyle{
				Color:    "#FF0000",
				FontSize: 18,
			},
		},
		BackgroundColor: "#333333",
		//ExpenseItems: expenseItems,

	}
	return res, nil

}

//func GetAllExpenseItemsForUser(userId int64)[]models.ExpenseItem{
//
//	items := []string{
//		"Movie expense",
//		"Trip expense",
//	}
//
//	var expenseList []models.ExpenseItem
//
//	//for _, item := range items {
//	//	expenseItem := models.ExpenseItem{
//	//		ItemName: models.UikitText{
//	//			Value: item,
//	//		} ,
//	//	}
//	//	expenseList = append(expenseList, expenseItem)
//	//}
//	return expenseList
//}

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
	}

	// 2. Wrap in the top-level envelope
	return models.ChatResponse{
		Type:      "chat_list",
		Chats:     chats,
	}
}

func intPtr(i int) *int { return &i }