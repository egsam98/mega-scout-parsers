package controllers

import (
	"github.com/egsam98/MegaScout/parsers"
	"github.com/egsam98/MegaScout/utils/errors"
	"github.com/gin-gonic/gin"
)

func LeagueDetailController(c *gin.Context) {
	url, exists := c.GetQuery("url")
	if !exists {
		c.Error(errors.NewClientError(400, "'url' is not provided"))
		return
	}

	league, err := parsers.LeagueDetail(url)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, league)
}
