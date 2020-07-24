package controllers

import (
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/parsers"
	"github.com/egsam98/MegaScout/utils/errors"
	"github.com/gin-gonic/gin"
)

type Zopa struct {
	models.Goal //swagger:allOf
	models.Card // swagger:allOf
	models.Penalty
	models.Substitution
}

// @Router /match_events [get]
// @Summary Список событий матча
// @Param match_url query string true "URL матча"
// @Produce json
// @Success 200 {array} models.AllMatchEventFields
// @Failure 400 {object} models.ErrorJSON
// @Failure 408
func MatchEventsController(c *gin.Context) {
	matchUrl := c.Query("match_url")
	if matchUrl == "" {
		c.Error(errors.NewClientError(400, "match_url is not provided"))
		return
	}

	data, err := parsers.MatchEvents(matchUrl)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}
