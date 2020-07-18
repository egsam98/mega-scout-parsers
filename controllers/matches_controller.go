package controllers

import (
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
)

// @Router /matches [get]
// @Summary Список матчей клуба
// @Param team_url query string true "URL клуба"
// @Produce json
// @Success 200 {array} models.Match
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func MatchesController(c *gin.Context) {
	teamUrl := c.Query("team_url")
	if teamUrl == "" {
		c.JSON(400, models.NewErrorJSON("team_url is not provided"))
		return
	}

	data, err := parsers.Matches(teamUrl)
	if err != nil {
		c.Status(408)
		return
	}

	c.JSON(200, data)
}
