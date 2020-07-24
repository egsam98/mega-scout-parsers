package main

import (
	. "github.com/egsam98/MegaScout/controllers"
	"github.com/egsam98/MegaScout/docs"
	"github.com/egsam98/MegaScout/middlewares"
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
	r.Use(middlewares.JSONErrorHandler())
	r.GET("/countries", CountriesController)
	r.GET("/seasons", SeasonsController)
	r.GET("/leagues", LeaguesController)
	r.GET("/team_compositions", TeamCompositionsController)
	r.GET("/player_detail", PlayerDetailController)
	r.GET("/team_detail", TeamDetailController)
	r.GET("/trophies", TrophiesController)
	r.GET("/matches", MatchesController)
	r.GET("/match_events", MatchEventsController)
	r.GET("/coaches", CoachesController)
	r.GET("/player_stats", PlayerStatsController)

	initSwagger(r)

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}
