package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/egsam98/MegaScout/utils/errors"
	"github.com/gin-gonic/gin"
)

// @Router /player_detail [get]
// @Summary Карточка игрока
// @Param url query string true "URL игрока"
// @Produce json
// @Success 200 {object} models.PlayerDetail
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func PlayerDetailController(c *gin.Context) {
	playerUrl := c.Query("url")
	if playerUrl == "" {
		c.Error(errors.NewClientError(400, "url is not provided"))
		return
	}
	data, err := parsers.PlayerDetail(playerUrl)
	if err != nil {
		c.Status(408)
		return
	}
	c.JSON(200, data)
}
