package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/message"
	"github.com/egsam98/MegaScout/utils/pointers"
	strings2 "github.com/egsam98/MegaScout/utils/strings"
	"strconv"
	"strings"
)

func PlayerStats(playerUrl string, seasonPeriod *int) ([]models.PlayerStats, error) {
	url := strings.ReplaceAll(playerUrl, "profil", "leistungsdatendetails")
	statsChan := make(chan message.Message)

	var seasonPeriods []int
	if seasonPeriod != nil {
		seasonPeriods = []int{*seasonPeriod}
	} else {
		for _, season := range seasons {
			seasonPeriods = append(seasonPeriods, season.Period)
		}
	}

	for _, period := range seasonPeriods {
		go processStats(url, period, statsChan)
	}

	stats := make([]models.PlayerStats, 0)
	for range seasonPeriods {
		msg := <-statsChan
		if msg.IsError() {
			return nil, msg.Error
		}
		stats = append(stats, msg.Data.([]models.PlayerStats)...)
	}
	return stats, nil
}

func processStats(url string, seasonPeriod int, statsChan chan<- message.Message) {
	doc, err := utils.FetchHtml(url + "?saison=" + strconv.Itoa(seasonPeriod) + "&plus=1")
	if err != nil {
		statsChan <- message.Error(err)
		return
	}

	stats := make([]models.PlayerStats, 0)
	boxes := doc.Find(".large-12.columns > .box")
	boxes.Slice(1, boxes.Length()).Each(func(_ int, box *goquery.Selection) {
		box.Find("tbody tr").Each(func(_ int, tr *goquery.Selection) {
			tds := tr.Find("td")
			matchId, err := strconv.Atoi(tds.Eq(-11).Find("a").AttrOr("id", ""))
			if err == nil {
				var subOn *int
				var subOff *int
				if result := tds.Eq(-2).Text(); result != "" {
					subOff = pointers.NewInt(parseGameMinutes(result))
				}
				if result := tds.Eq(-3).Text(); result != "" {
					subOn = pointers.NewInt(parseGameMinutes(result))
				}

				stats = append(stats, models.PlayerStats{
					MatchId:           matchId,
					SeasonPeriod:      seasonPeriod,
					Goals:             strings2.ToInt(tds.Eq(-9).Text(), true),
					Assists:           strings2.ToInt(tds.Eq(-8).Text(), true),
					OwnGoals:          strings2.ToInt(tds.Eq(-7).Text(), true),
					YellowCards:       strings2.ToInt(tds.Eq(-6).Text(), true),
					SecondYellowCards: strings2.ToInt(tds.Eq(-5).Text(), true),
					RedCards:          strings2.ToInt(tds.Eq(-4).Text(), true),
					SubstitutionOn:    subOn,
					SubstitutionOff:   subOff,
					MinutesPlayed:     parseGameMinutes(tds.Eq(-1).Text()),
				})
			}
		})
	})
	statsChan <- message.Ok(stats)
}

func parseGameMinutes(minutes string) int {
	return strings2.ToInt(strings.Trim(minutes, "'\u00a0"), false)
}
