package controllers

import (
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
)

// @Router /coaches [get]
// @Summary 2 тренера матча (в доработке)
// @Param match_url query string true "URL матча"
// @Produce json
// @Success 200 {array} models.Coach
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func CoachesController(c *gin.Context) {
	matchUrl := c.Query("match_url")
	if matchUrl == "" {
		c.JSON(400, models.NewErrorJSON("match_url is not provided"))
		return
	}

	data, err := parsers.Coaches(matchUrl)
	if err != nil {
		c.Status(408)
		return
	}

	c.JSON(200, data)
}
