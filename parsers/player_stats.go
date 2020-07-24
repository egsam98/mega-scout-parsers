package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/message"
	"github.com/egsam98/MegaScout/utils/pointers"
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

	defer func() {
		obj := recover()
		if err, ok := obj.(error); ok {
			statsChan <- message.Error(err)
		}
	}()

	stats := make([]models.PlayerStats, 0)
	boxes := doc.Find(".large-12.columns > .box")
	boxes.Slice(1, boxes.Length()).Each(func(_ int, box *goquery.Selection) {
		box.Find("tbody tr").Each(func(_ int, tr *goquery.Selection) {
			tds := tr.Find("td")
			matchId, err := strconv.Atoi(tds.Eq(-11).Find("a").AttrOr("id", ""))
			if err == nil {
				minutesPlayed, err := parseGameMinutes(tds.Eq(-1).Text())
				if err != nil {
					panic(err)
				}

				var subOn *int
				var subOff *int
				if result, err := parseGameMinutes(tds.Eq(-2).Text()); err == nil {
					subOff = pointers.NewInt(result)
				}
				if result, err := parseGameMinutes(tds.Eq(-3).Text()); err == nil {
					subOn = pointers.NewInt(result)
				}

				goals, _ := strconv.Atoi(tds.Eq(-9).Text())
				assists, _ := strconv.Atoi(tds.Eq(-8).Text())
				ownGoals, _ := strconv.Atoi(tds.Eq(-7).Text())
				yellowCards, _ := strconv.Atoi(tds.Eq(-6).Text())
				secondYellowCards, _ := strconv.Atoi(tds.Eq(-5).Text())
				redCards, _ := strconv.Atoi(tds.Eq(-4).Text())

				stats = append(stats, models.PlayerStats{
					MatchId:           matchId,
					SeasonPeriod:      seasonPeriod,
					Goals:             goals,
					Assists:           assists,
					OwnGoals:          ownGoals,
					YellowCards:       yellowCards,
					SecondYellowCards: secondYellowCards,
					RedCards:          redCards,
					SubstitutionOn:    subOn,
					SubstitutionOff:   subOff,
					MinutesPlayed:     minutesPlayed,
				})
			}
		})
	})
	statsChan <- message.Ok(stats)
}

func parseGameMinutes(minutes string) (int, error) {
	return strconv.Atoi(strings.Trim(minutes, "'\u00a0"))
}
