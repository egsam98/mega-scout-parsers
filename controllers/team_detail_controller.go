package controllers

import (
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
)

// @Router /team_detail [get]
// @Summary Карточка клуба
// @Param url query string true "URL клуба"
// @Produce json
// @Success 200 {object} models.TeamDetail
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func TeamDetailController(c *gin.Context) {
	teamUrl := c.Query("url")
	if teamUrl == "" {
		c.JSON(400, models.NewErrorJSON("url is not provided"))
		return
	}

	data, err := parsers.TeamDetail(teamUrl)
	if err != nil {
		c.Status(408)
		return
	}

	c.JSON(200, data)
}
