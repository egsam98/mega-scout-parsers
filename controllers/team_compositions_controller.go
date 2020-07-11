package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
)

func TeamCompositionsController(c *gin.Context) {
	leagueUrl := c.Query("league_url")
	if len(leagueUrl) == 0 {
		c.JSON(400, gin.H{
			"error": "league_url is not provided",
		})
		return
	}
	data, err := parsers.TeamCompositions(leagueUrl)
	if err != nil {
		c.Status(408)
		return
	}
	c.JSON(200, data)
}
