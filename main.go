package main

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load()
	r := gin.Default()
	r.GET("/countries", parseCountries)
	r.GET("/seasons", parseSeasons)
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}

func parseCountries(c *gin.Context) {
	data, err := parsers.Countries()
	if err != nil {
		panic(err)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(200, data)
}

func parseSeasons(c *gin.Context) {
	seasons := parsers.Seasons()
	c.Header("Content-Type", "application/json")
	c.JSON(200, seasons)
}
