package players

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db"
	handlerErrors "github.com/mellena1/boston-archery-api/handlers/errors"
	"github.com/mellena1/boston-archery-api/ptr"
)

var failedToUpdatePlayerError = handlerErrors.Error{
	Msg: "failed to update player",
}

// swagger:parameters putPlayer
type PutPlayerInput struct {
	// in:path
	ID uuid.UUID
	// in:body
	Body struct {
		// required: true
		FirstName string `json:"firstName" binding:"required,ascii"`
		// required: true
		LastName string `json:"lastName" binding:"required"`
	}
}

// swagger:model PutPlayerResp
type PutPlayerResp struct {
	Data Player `json:"data"`
}

// swagger:route PUT /player/{ID} player putPlayer
//
// Update a player.
//
// responses:
//
//	200: body:PutPlayerResp
//	400: body:Error
//	500: body:Error
func (a *API) PutPlayer(c *gin.Context) {
	var input PutPlayerInput
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

	player, err := a.db.UpdatePlayer(c.Request.Context(), input.ID, db.UpdatePlayerInput{
		FirstName: ptr.Ptr(input.Body.FirstName),
		LastName:  ptr.Ptr(input.Body.LastName),
	})
	if err != nil {
		switch {
		case errors.Is(err, db.ErrItemAlreadyExists):
			c.AbortWithStatusJSON(http.StatusConflict, handlerErrors.AlreadyExistsError)
		default:
			a.logger.ErrorContext(c, "failed to add player to db", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, failedToUpdatePlayerError)
		}
		return
	}

	c.JSON(http.StatusOK, PostPlayerResp{
		Data: playerFromModel(*player),
	})
}
