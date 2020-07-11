package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
)

func SeasonsController(c *gin.Context) {
	seasons := parsers.Seasons()
	c.Header("Content-Type", "application/json")
	c.JSON(200, seasons)
}
