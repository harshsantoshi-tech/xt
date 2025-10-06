package controller

import (
	"expense-tracker/handlers"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HomePageController(c echo.Context) error {

	res, err := handlers.HomePageHandler(1)
	if err != nil {
		log.Println("Error giving home page data ", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}
