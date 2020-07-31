package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/egsam98/MegaScout/utils/errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Router /team_compositions [get]
// @Summary Список команд и игроков в них по соревнованию (лига & сезон)
// @Param league_url query string true "URL лиги"
// @Param season_period query int false "период сезона (год). По умолч. последний период"
// @Produce json
// @Success 200 {array} models.Team
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func TeamCompositionsController(c *gin.Context) {
	leagueUrl, exists := c.GetQuery("league_url")
	if !exists {
		c.Error(errors.NewClientError(400, "league_url is not provided"))
		return
	}

	var seasonPeriod int
	if result, exists := c.GetQuery("season_period"); exists {
		var err error
		seasonPeriod, err = strconv.Atoi(result)
		if err != nil {
			c.Error(err)
			return
		}
	} else {
		seasonPeriod = parsers.LatestSeason().Period
	}

	data, err := parsers.TeamCompositions(leagueUrl, seasonPeriod)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, data)
}
