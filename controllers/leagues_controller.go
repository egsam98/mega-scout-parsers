package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
	"strconv"
)

func LeaguesController(c *gin.Context) {
	countryId, err := strconv.Atoi(c.Query("country"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid country. Must be integer.",
		})
		return
	}
	seasonPeriod, err := strconv.Atoi(c.Query("season_period"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid season_period. Must be integer.",
		})
		return
	}
	leagues, err := parsers.Leagues(countryId, seasonPeriod)
	if err != nil {
		c.Status(408)
		return
	}
	c.JSON(200, leagues)
}
