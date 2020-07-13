package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
)

// @Router /countries [get]
// @Summary Список стран
// @Produce json
// @Success 200 {array} models.Country
func CountriesController(c *gin.Context) {
	data, err := parsers.Countries()
	if err != nil {
		panic(err)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(200, data.Slice())
}
