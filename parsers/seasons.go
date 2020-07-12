package parsers

import (
	"github.com/egsam98/MegaScout/models"
	"time"
)

func Seasons() []models.Season {
	today := time.Now()
	seasons := make([]models.Season, 0, today.Year()+1-1900)
	for i := 1900; i <= today.Year(); i++ {
		seasonStart := time.Date(i, 7, 1, 0, 0, 0, 0, time.UTC)
		if seasonStart.After(today) {
			continue
		}
		seasons = append(seasons, models.Season{
			Period:      i,
			SeasonStart: seasonStart.Format("02-01-2006"),
			SeasonEnd: time.Date(i+1, 6, 30, 0, 0, 0, 0, time.UTC).
				Format("02-01-2006"),
		})
	}
	return seasons
}
