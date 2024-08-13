package teams

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mellena1/boston-archery-api/db"
	handlerErrors "github.com/mellena1/boston-archery-api/handlers/errors"
	"github.com/mellena1/boston-archery-api/model"
)

var failedToAddTeamError = handlerErrors.Error{
	Msg: "failed to add team",
}

// swagger:parameters postTeam
type PostTeamInput struct {
	// in:body
	Body struct {
		// required: true
		Name string `json:"name" binding:"required,ascii"`
		// required: true
		TeamColors []string `json:"teamColors" binding:"required"`
	}
}

// swagger:model PostTeamResp
type PostTeamResp struct {
	Data Team `json:"data"`
}

// swagger:route POST /team team postTeam
//
// Add a new team.
//
// responses:
//
//	200: body:PostTeamResp
//	400: body:Error
//	500: body:Error
func (a *API) PostTeam(c *gin.Context) {
	var input PostTeamInput
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

	team, err := a.db.AddTeam(c.Request.Context(), model.Team{
		ID:         uuid.New(),
		Name:       input.Body.Name,
		TeamColors: input.Body.TeamColors,
	})
	if err != nil {
		switch {
		case errors.Is(err, db.ErrItemAlreadyExists):
			c.AbortWithStatusJSON(http.StatusConflict, handlerErrors.AlreadyExistsError)
		default:
			a.logger.ErrorContext(c, "failed to add team to db", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, failedToAddTeamError)
		}
		return
	}

	c.JSON(http.StatusOK, PostTeamResp{
		Data: teamFromModel(*team),
	})
}

var validHexColorRegex regexp.Regexp = *regexp.MustCompile("#[0-9a-fA-F]{6}")

func validateTeamColors(colors []string) error {
	for _, color := range colors {
		if !validHexColorRegex.MatchString(color) {
			return fmt.Errorf("invalid color. should be hex: %s", color)
		}
	}

	return nil
}
