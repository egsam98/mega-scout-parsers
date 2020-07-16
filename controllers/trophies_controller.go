package controllers

import (
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
)

// @Router /trophies [get]
// @Summary Трофеи игрока/тренера
// @Param person_url query string true "URL игрока/тренера"
// @Produce json
// @Success 200 {array} models.Trophy
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func TrophiesController(c *gin.Context) {
	personUrl := c.Query("person_url")
	if personUrl == "" {
		c.JSON(400, models.NewErrorJSON("person_url (player or coach) is not provided"))
		return
	}

	data, err := parsers.Trophies(personUrl)
	if err != nil {
		c.Status(408)
		return
	}

	c.JSON(200, data)
}
