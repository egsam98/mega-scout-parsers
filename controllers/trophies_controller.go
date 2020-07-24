package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/egsam98/MegaScout/utils/errors"
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
		c.Error(errors.NewClientError(400, "person_url (player or coach) is not provided"))
		return
	}

	data, err := parsers.Trophies(personUrl)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}
