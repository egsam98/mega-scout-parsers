package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/egsam98/MegaScout/utils/errors"
	"github.com/gin-gonic/gin"
)

// @Router /team_compositions [get]
// @Summary Список команд и игроков в них по соревнованию (лига & сезон)
// @Param league_url query string true "URL соревнования"
// @Produce json
// @Success 200 {array} models.Team
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func TeamCompositionsController(c *gin.Context) {
	leagueUrl := c.Query("league_url")
	if len(leagueUrl) == 0 {
		c.Error(errors.NewClientError(400, "league_url is not provided"))
		return
	}
	data, err := parsers.TeamCompositions(leagueUrl)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, data)
}
