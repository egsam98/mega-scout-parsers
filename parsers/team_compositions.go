package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"strconv"
)

func TeamCompositions(leagueUrl string, seasonPeriod int) (teams []models.Team, _ error) {
	doc, err := utils.FetchHtml(leagueUrl + "/saison_id/" + strconv.Itoa(seasonPeriod))
	if err != nil {
		return nil, err
	}
	var innerError error
	count := 0
	teamChan := make(chan models.Team)
	doc.Find("#yw1 tbody > tr > td:nth-child(2) > a:nth-child(1)").Each(func(i int, a *goquery.Selection) {
		if len(a.Text()) == 0 {
			return
		}
		count++
		go func(a *goquery.Selection) {
			url, _ := a.Attr("href")
			players, err := processPlayers(BaseUrl + url)
			if err != nil {
				innerError = err
				return
			}

			id, _ := strconv.Atoi(a.AttrOr("id", ""))
			teamChan <- models.Team{
				Id:      id,
				Url:     BaseUrl + url,
				Title:   a.Text(),
				Players: players,
			}
		}(a)
	})
	for i := 0; i < count; i++ {
		teams = append(teams, <-teamChan)
	}
	if innerError != nil {
		return nil, innerError
	}
	return teams, nil
}

func processPlayers(clubUrl string) (players []models.Player, _ error) {
	doc, err := utils.FetchHtml(clubUrl)
	if err != nil {
		return nil, err
	}
	doc.Find("#yw1 tbody > tr").Each(func(i int, tr *goquery.Selection) {
		if _, exists := tr.Attr("class"); !exists {
			return
		}
		tds := tr.Children()
		a := tds.Find(".hauptlink a").First()
		id, _ := strconv.Atoi(a.AttrOr("id", ""))
		players = append(players, models.Player{
			Id:  id,
			Url: BaseUrl + a.AttrOr("href", ""),
		})
	})
	return players, nil
}
