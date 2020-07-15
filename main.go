package main

import (
	. "github.com/egsam98/MegaScout/controllers"
	"github.com/egsam98/MegaScout/docs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"runtime"
)

func initSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "MegaScout Parsers API"
	docs.SwaggerInfo.Version = "1.0"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	_ = godotenv.Load()

	r := gin.Default()
	r.GET("/countries", CountriesController)
	r.GET("/seasons", SeasonsController)
	r.GET("/leagues", LeaguesController)
	r.GET("/team_compositions", TeamCompositionsController)
	r.GET("/player_detail", PlayerDetailController)
	r.GET("/team_detail", TeamDetailController)

	initSwagger(r)

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}
