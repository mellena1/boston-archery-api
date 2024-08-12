package seasons

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db"
	"github.com/mellena1/boston-archery-api/handlers"
	handlerErrors "github.com/mellena1/boston-archery-api/handlers/errors"
	"github.com/mellena1/boston-archery-api/model"
	"github.com/mellena1/boston-archery-api/slices"
)

var failedToUpdateSeasonError = handlerErrors.Error{
	Msg: "failed to update season",
}

// swagger:parameters putSeason
type PutSeasonInput struct {
	// in:path
	ID uuid.UUID
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

// swagger:model PutSeasonResp
type PutSeasonResp struct {
	Data Season `json:"data"`
}

// swagger:route PUT /season/{ID} season putSeason
//
// Update a season.
//
// responses:
//
//	200: body:PutSeasonResp
//	400: body:Error
//	500: body:Error
func (a *API) PutSeason(c *gin.Context) {
	var input PutSeasonInput
	var err error

	strId := c.Param("id")
	input.ID, err = uuid.Parse(strId)
	if err != nil {
		a.logger.WarnContext(c, "invalid ID", slog.String("id", strId), slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, handlerErrors.BadRequestError)
		return
	}

	if err := c.ShouldBindJSON(&input.Body); err != nil {
		a.logger.WarnContext(c, "failed to bind json", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, handlerErrors.BadRequestError)
		return
	}

	if err := validateByeWeeks(input.Body.StartDate.ToTime(), input.Body.EndDate.ToTime(), input.Body.ByeWeeks); err != nil {
		a.logger.WarnContext(c, "failed to validate bye weeks", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, handlerErrors.BadRequestError)
		return
	}

	season, err := a.db.UpdateSeason(c.Request.Context(), model.Season{
		ID:        input.ID,
		Name:      input.Body.Name,
		StartDate: input.Body.StartDate.ToTime(),
		EndDate:   input.Body.EndDate.ToTime(),
		ByeWeeks: slices.Map(input.Body.ByeWeeks, func(v handlers.Date) time.Time {
			return v.ToTime()
		}),
	})
	if err != nil {
		switch {
		case errors.Is(err, db.ErrItemNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, handlerErrors.NotFoundError)
		default:
			a.logger.ErrorContext(c, "failed to update season to db", slog.String("id", input.ID.String()), slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, failedToUpdateSeasonError)
		}
		return
	}

	c.JSON(http.StatusOK, PostSeasonResp{
		Data: seasonFromModel(*season),
	})
}
