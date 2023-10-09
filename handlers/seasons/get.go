package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/handlers/errors"
)

var failedToFetchSeasonsError = errors.Error{
	Msg: "failed to fetch seasons",
}

func (a *API) GetSeasons(c *gin.Context) {
	seasons, err := a.db.GetAllSeasons(c.Request.Context())
	if err != nil {
		a.logger.Error("failed to get seasons", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, failedToFetchSeasonsError)
		return
	}

	c.JSON(http.StatusOK, seasons)
}
