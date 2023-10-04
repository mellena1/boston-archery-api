package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/handlers"
)

var failedToFetchSeasonsError = handlers.Error{
	Msg: "failed to fetch seasons",
}

func GetSeasons(c *gin.Context) {
	seasons, err := database.GetAllSeasons(c.Request.Context())
	if err != nil {
		logger.Error("failed to get seasons", "error", err)
		c.JSON(http.StatusInternalServerError, failedToFetchSeasonsError)
		return
	}

	c.JSON(http.StatusOK, seasons)
}
