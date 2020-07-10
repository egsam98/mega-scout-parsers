package main

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	_ = godotenv.Load()
	r := gin.Default()
	r.GET("/countries", parseCountries)
	r.GET("/seasons", parseSeasons)
	r.GET("/leagues", parseLeagues)
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

func parseLeagues(c *gin.Context) {
	countryId, err := strconv.Atoi(c.Query("country"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid country. Must be integer.",
		})
		return
	}
	seasonPeriod, err := strconv.Atoi(c.Query("season_period"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid season_period. Must be integer.",
		})
		return
	}
	leagues, err := parsers.Leagues(countryId, seasonPeriod)
	if err != nil {
		c.JSON(408, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, leagues)
}
