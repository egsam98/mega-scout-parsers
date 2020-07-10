package main

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	r := gin.Default()
	r.GET("/countries", func(c *gin.Context) {
		data, err := parsers.Countries()
		if err != nil {
			panic(err)
		}
		c.Header("Content-Type", "application/json")
		c.JSON(200, data)
	})
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}
