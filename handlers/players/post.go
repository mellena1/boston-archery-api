package players

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/db"
	handlerErrors "github.com/mellena1/boston-archery-api/handlers/errors"
)

var failedToAddPlayerError = handlerErrors.Error{
	Msg: "failed to add player",
}

// swagger:parameters postPlayer
type PostPlayerInput struct {
	// in:body
	Body struct {
		// required: true
		FirstName string `json:"firstName" binding:"required,ascii"`
		// required: true
		LastName string `json:"lastName" binding:"required"`
	}
}

// swagger:model PostPlayerResp
type PostPlayerResp struct {
	Data Player `json:"data"`
}

// swagger:route POST /player player postPlayer
//
// Add a new player.
//
// responses:
//
//	200: body:PostPlayerResp
//	400: body:Error
//	500: body:Error
func (a *API) PostPlayer(c *gin.Context) {
	var input PostPlayerInput
	if err := c.ShouldBindJSON(&input.Body); err != nil {
		a.logger.ErrorContext(c, "failed to bind json", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, handlerErrors.BadRequestError)
		return
	}

	player, err := a.db.AddPlayer(c.Request.Context(), db.PlayerInput{
		FirstName: input.Body.FirstName,
		LastName:  input.Body.LastName,
	})
	if err != nil {
		switch {
		case errors.Is(err, db.ErrItemAlreadyExists):
			c.AbortWithStatusJSON(http.StatusConflict, handlerErrors.AlreadyExistsError)
		default:
			a.logger.ErrorContext(c, "failed to add player to db", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, failedToAddPlayerError)
		}
		return
	}

	c.JSON(http.StatusOK, PostPlayerResp{
		Data: playerFromModel(*player),
	})
}
