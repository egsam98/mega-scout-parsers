package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
)

// @Router /seasons [get]
// @Summary Список сезонов с 1900г. по текущий
// @Produce json
// @Success 200 {array} models.Season
func SeasonsController(c *gin.Context) {
	seasons := parsers.Seasons()
	c.Header("Content-Type", "application/json")
	c.JSON(200, seasons)
}
