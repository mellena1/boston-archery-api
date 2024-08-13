package teams

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db"
	handlerErrors "github.com/mellena1/boston-archery-api/handlers/errors"
)

var failedToUpdateTeamError = handlerErrors.Error{
	Msg: "failed to update team",
}

// swagger:parameters putTeam
type PutTeamInput struct {
	// in:path
	ID uuid.UUID
	// in:body
	Body struct {
		// required: true
		Name string `json:"name" binding:"required,ascii"`
		// required: true
		TeamColors []string `json:"teamColors" binding:"required"`
	}
}

// swagger:model PutTeamResp
type PutTeamResp struct {
	Data Team `json:"data"`
}

// swagger:route PUT /team/{ID} team putTeam
//
// Update a team.
//
// responses:
//
//	200: body:PutTeamResp
//	400: body:Error
//	500: body:Error
func (a *API) PutTeam(c *gin.Context) {
	var input PutTeamInput
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

	if err := validateTeamColors(input.Body.TeamColors); err != nil {
		a.logger.WarnContext(c, "failed to validate team colors", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, handlerErrors.BadRequestError)
		return
	}

	team, err := a.db.UpdateTeam(c.Request.Context(), input.ID, db.UpdateTeamInput{
		Name:       &input.Body.Name,
		TeamColors: &input.Body.TeamColors,
	})
	if err != nil {
		switch {
		case errors.Is(err, db.ErrItemNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, handlerErrors.NotFoundError)
		default:
			a.logger.ErrorContext(c, "failed to update team to db", slog.String("id", input.ID.String()), slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, failedToUpdateTeamError)
		}
		return
	}

	c.JSON(http.StatusOK, PutTeamResp{
		Data: teamFromModel(*team),
	})
}
