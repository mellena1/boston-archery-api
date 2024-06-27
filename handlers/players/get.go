package players

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handlerErrors "github.com/mellena1/boston-archery-api/handlers/errors"
	"github.com/mellena1/boston-archery-api/model"
	"github.com/mellena1/boston-archery-api/slices"
)

var failedToFetchPlayerError = handlerErrors.Error{
	Msg: "failed to fetch player(s)",
}

// swagger:model Player
type Player struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func playerFromModel(player model.Player) Player {
	return Player{
		ID:        player.ID.String(),
		FirstName: player.FirstName,
		LastName:  player.LastName,
	}
}

// swagger:parameters getPlayer
type GetPlayerInput struct {
	// in:path
	ID uuid.UUID
}

// swagger:model GetPlayerResp
type GetPlayerResp struct {
	Data Player `json:"data"`
}

// swagger:route GET /player/{ID} player getPlayer
//
// Get a player by ID.
//
// responses:
//
//	200: body:GetPlayerResp
//	400: body:Error
//	500: body:Error
func (a *API) GetPlayer(c *gin.Context) {
	var input GetPlayerInput
	var err error

	strId := c.Param("id")
	input.ID, err = uuid.Parse(strId)
	if err != nil {
		a.logger.WarnContext(c, "failed to parse id", slog.String("id", strId), slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, handlerErrors.BadRequestError)
		return
	}

	player, err := a.db.GetPlayer(c, input.ID)
	if err != nil {
		a.logger.ErrorContext(c, "failed to get player from db", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, failedToFetchPlayerError)
		return
	}

	c.JSON(http.StatusOK, GetPlayerResp{
		Data: playerFromModel(*player),
	})
}

// swagger:model GetPlayersResp
type GetPlayersResp struct {
	Data []Player `json:"data"`
}

// swagger:route GET /players player getPlayers
//
// Get all players.
//
// responses:
//
//	200: body:GetPlayersResp
//	500: body:Error
func (a *API) GetPlayers(c *gin.Context) {
	players, err := a.db.GetAllPlayers(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, failedToFetchPlayerError)
		return
	}

	c.JSON(http.StatusOK, GetPlayersResp{
		Data: slices.Map(players, func(p model.Player) Player {
			return playerFromModel(p)
		}),
	})
}
