package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/handlers"
	"github.com/mellena1/boston-archery-api/slices"
)

var failedToAddSeasonError = handlers.Error{
	Msg: "failed to add season",
}

type PostSeasonInput struct {
	Name      string          `json:"name" binding:"required"`
	StartDate handlers.Date   `json:"startDate" binding:"required"`
	EndDate   handlers.Date   `json:"endDate" binding:"required,gtfield=StartDate"`
	ByeWeeks  []handlers.Date `json:"byeWeeks"`
}

func PostSeason(c *gin.Context) {
	var input PostSeasonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Debug("bad request", "error", err)
		c.JSON(http.StatusBadRequest, handlers.BadRequestError)
		return
	}

	err := database.AddSeason(c.Request.Context(), db.SeasonInput{
		Name:      input.Name,
		StartDate: input.StartDate.ToTime(),
		EndDate:   input.EndDate.ToTime(),
		ByeWeeks: slices.Map(input.ByeWeeks, func(v handlers.Date) time.Time {
			return v.ToTime()
		}),
	})
	if err != nil {
		logger.Error("failed to add season to db", "error", err)
		c.JSON(http.StatusInternalServerError, failedToAddSeasonError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
