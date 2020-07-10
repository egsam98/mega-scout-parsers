package parsers

import (
	"github.com/gin-gonic/gin"
	"time"
)

func Seasons() []gin.H {
	today := time.Now()
	seasons := make([]gin.H, 0, today.Year()+1-1900)
	for i := 1900; i <= today.Year(); i++ {
		seasonStart := time.Date(i, 7, 1, 0, 0, 0, 0, time.UTC)
		if seasonStart.After(today) {
			continue
		}
		seasons = append(seasons, gin.H{
			"period":       i,
			"season_start": seasonStart,
			"season_end":   time.Date(i+1, 6, 30, 0, 0, 0, 0, time.UTC),
		})
	}
	return seasons
}
