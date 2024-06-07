package seasons

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/handlers"
	handlerErrors "github.com/mellena1/boston-archery-api/handlers/errors"
	"github.com/mellena1/boston-archery-api/slices"
)

var failedToAddSeasonError = handlerErrors.Error{
	Msg: "failed to add season",
}

// swagger:parameters postSeason
type PostSeasonInput struct {
	// in:body
	Body struct {
		// required: true
		Name string `json:"name" binding:"required,min=5,ascii"`
		// required: true
		StartDate handlers.Date `json:"startDate" binding:"required"`
		// required: true
		EndDate  handlers.Date   `json:"endDate" binding:"required,gtfield=StartDate"`
		ByeWeeks []handlers.Date `json:"byeWeeks" binding:"unique"`
	}
}

// swagger:model PostSeasonResp
type PostSeasonResp struct {
	Data Season `json:"data"`
}

// swagger:route POST /seasons seasons postSeason
//
// Add a new season.
//
// responses:
//
//	200: body:PostSeasonResp
//	400: body:Error
//	500: body:Error
func (a *API) PostSeason(c *gin.Context) {
	var input PostSeasonInput
	if err := c.ShouldBindJSON(&input.Body); err != nil {
		a.logger.Error("failed to bind json", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, handlerErrors.BadRequestError)
		return
	}

	if err := validateByeWeeks(input.Body.StartDate.ToTime(), input.Body.EndDate.ToTime(), input.Body.ByeWeeks); err != nil {
		a.logger.Error("failed to validate bye weeks", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, handlerErrors.BadRequestError)
		return
	}

	season, err := a.db.AddSeason(c.Request.Context(), db.SeasonInput{
		Name:      input.Body.Name,
		StartDate: input.Body.StartDate.ToTime(),
		EndDate:   input.Body.EndDate.ToTime(),
		ByeWeeks: slices.Map(input.Body.ByeWeeks, func(v handlers.Date) time.Time {
			return v.ToTime()
		}),
	})
	if err != nil {
		switch {
		case errors.Is(err, db.ErrItemAlreadyExists):
			c.AbortWithStatusJSON(http.StatusConflict, handlerErrors.AlreadyExistsError)
		default:
			a.logger.Error("failed to add season to db", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, failedToAddSeasonError)
		}
		return
	}

	c.JSON(http.StatusOK, PostSeasonResp{
		Data: seasonFromModel(*season),
	})
}

func validateByeWeeks(startDate, endDate time.Time, byeWeeks []handlers.Date) error {
	for _, bye := range byeWeeks {
		byeAsTime := bye.ToTime()
		if startDate.After(byeAsTime) || endDate.Before(byeAsTime) {
			return fmt.Errorf("bye week is not in season")
		}
	}

	return nil
}
