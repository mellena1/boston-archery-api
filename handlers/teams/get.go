package teams

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	handlerErrors "github.com/mellena1/boston-archery-api/handlers/errors"
	"github.com/mellena1/boston-archery-api/model"
	"github.com/mellena1/boston-archery-api/slices"
)

var failedToFetchTeamsError = handlerErrors.Error{
	Msg: "failed to fetch teams",
}

// swagger:parameters getTeams
type GetTeamsInput struct{}

// swagger:model GetTeamsResp
type GetTeamsResp struct {
	Data []Team `json:"data"`
}

// swagger:model Team
type Team struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	TeamColors []string `json:"teamColors"`
}

func teamFromModel(t model.Team) Team {
	return Team{
		ID:         t.ID.String(),
		Name:       t.Name,
		TeamColors: t.TeamColors,
	}
}

// swagger:route GET /teams team getTeams
// Get all teams.
//
// responses:
//
//	200: body:GetTeamsResp
//	500: Error
func (a *API) GetTeams(c *gin.Context) {
	var input GetTeamsInput
	err := c.BindQuery(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, handlerErrors.BadRequestError)
		return
	}

	a.getAllTeams(c)
}

func (a *API) getAllTeams(c *gin.Context) {
	teams, err := a.db.GetAllTeams(c.Request.Context())
	if err != nil {
		a.logger.ErrorContext(c, "failed to get teams", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, failedToFetchTeamsError)
		return
	}

	c.JSON(http.StatusOK, GetTeamsResp{
		Data: slices.Map(teams, func(t model.Team) Team {
			return teamFromModel(t)
		}),
	})
}
