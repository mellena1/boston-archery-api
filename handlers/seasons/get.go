package seasons

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mellena1/boston-archery-api/db"
	handlerErrors "github.com/mellena1/boston-archery-api/handlers/errors"
	"github.com/mellena1/boston-archery-api/model"
	"github.com/mellena1/boston-archery-api/slices"
)

var failedToFetchSeasonsError = handlerErrors.Error{
	Msg: "failed to fetch seasons",
}

// swagger:parameters getSeasons
type GetSeasonsInput struct {
	// in:query
	Name *string `form:"name"`
}

// swagger:model GetSeasonsResp
type GetSeasonsResp struct {
	Data []Season `json:"data"`
}

// swagger:model Season
type Season struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	ByeWeeks  []string `json:"byeWeeks,omitempty"`
}

func seasonFromModel(s model.Season) Season {
	return Season{
		ID:        s.ID.String(),
		Name:      s.Name,
		StartDate: s.StartDate.Format("2006-01-02"),
		EndDate:   s.EndDate.Format("2006-01-02"),
		ByeWeeks: slices.Map(s.ByeWeeks, func(t time.Time) string {
			return t.Format("2006-01-02")
		}),
	}
}

// swagger:route GET /seasons season getSeasons
// Get all seasons.
//
// responses:
//
//	200: body:GetSeasonsResp
//	500: Error
func (a *API) GetSeasons(c *gin.Context) {
	var input GetSeasonsInput
	err := c.BindQuery(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, handlerErrors.BadRequestError)
		return
	}

	if input.Name != nil {
		a.getSeasonByName(c, input)
		return
	}

	a.getAllSeasons(c)
}

func (a *API) getAllSeasons(c *gin.Context) {
	seasons, err := a.db.GetAllSeasons(c.Request.Context())
	if err != nil {
		a.logger.ErrorContext(c, "failed to get seasons", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, failedToFetchSeasonsError)
		return
	}

	c.JSON(http.StatusOK, GetSeasonsResp{
		Data: slices.Map(seasons, func(s model.Season) Season {
			return seasonFromModel(s)
		}),
	})
}

func (a *API) getSeasonByName(c *gin.Context, input GetSeasonsInput) {
	season, err := a.db.GetSeasonByName(c.Request.Context(), *input.Name)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrItemNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, handlerErrors.NotFoundError)
		default:
			a.logger.ErrorContext(c, "failed to get seasons", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, failedToFetchSeasonsError)
		}
		return
	}

	c.JSON(http.StatusOK, GetSeasonsResp{
		Data: []Season{seasonFromModel(*season)},
	})
}
