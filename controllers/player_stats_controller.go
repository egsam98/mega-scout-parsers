package controllers

import (
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Router /player_stats [get]
// @Summary Статистика игрока за каждый матч
// @Param player_url query string true "URL игрока"
// @Param season_period query string false "период сезона (год)"
// @Produce json
// @Success 200 {array} models.PlayerStats
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func PlayerStatsController(c *gin.Context) {
	playerUrl := c.Query("player_url")
	if playerUrl == "" {
		c.JSON(400, "player_url is not provided")
		return
	}

	var seasonPeriod *int
	seasonPeriodStr := c.Query("season_period")
	if seasonPeriodStr != "" {
		_seasonPeriod, err := strconv.Atoi(seasonPeriodStr)
		if err != nil {
			c.JSON(400, models.NewErrorJSON("season_period must be year"))
			return
		}
		seasonPeriod = &_seasonPeriod
	}

	data, err := parsers.PlayerStats(playerUrl, seasonPeriod)
	if err != nil {
		c.Status(408)
		return
	}
	c.JSON(200, data)
}
