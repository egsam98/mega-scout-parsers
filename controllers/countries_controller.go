package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
)

func CountriesController(c *gin.Context) {
	data, err := parsers.Countries()
	if err != nil {
		panic(err)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(200, data)
}
