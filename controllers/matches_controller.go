package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/egsam98/MegaScout/utils/errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Router /matches [get]
// @Summary Список матчей клуба
// @Param team_url query string true "URL клуба"
// @Param season_period query int false "период сезона (год)"
// @Produce json
// @Success 200 {array} models.Match
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func MatchesController(c *gin.Context) {
	teamUrl, ok := c.GetQuery("team_url")
	if !ok {
		c.Error(errors.NewClientError(400, "team_url is not provided"))
		return
	}

	var seasonPeriod *int
	if result, ok := c.GetQuery("season_period"); ok {
		result, err := strconv.Atoi(result)
		if err != nil {
			c.Error(errors.NewClientError(400, err.Error()))
			return
		}
		seasonPeriod = &result
	}

	data, err := parsers.Matches(teamUrl, seasonPeriod)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}
