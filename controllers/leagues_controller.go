package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/egsam98/MegaScout/utils/errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Router /leagues [get]
// @Summary Список лиг страны за определенный сезон
// @Param country query int true "ID страны (напр. 141)"
// @Param season_period query int true "Период сезона (напр. 2019)"
// @Produce json
// @Success 200 {array} models.League
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func LeaguesController(c *gin.Context) {
	countryId, err := strconv.Atoi(c.Query("country"))
	if err != nil {
		c.Error(errors.NewClientError(400, "Invalid country. Must be integer."))
		return
	}
	seasonPeriod, err := strconv.Atoi(c.Query("season_period"))
	if err != nil {
		c.Error(errors.NewClientError(400, "Invalid season_period. Must be integer."))
		return
	}
	leagues, err := parsers.Leagues(countryId, seasonPeriod)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, leagues)
}
