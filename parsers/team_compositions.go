package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func TeamCompositions(leagueUrl string) (teams []gin.H, _ error) {
	res, err := http.Get(leagueUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	var innerError error
	count := 0
	teamChan := make(chan gin.H)
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

			teamChan <- gin.H{
				"club_id": id,
				"url":     BaseUrl + url,
				"title":   a.Text(),
				"players": players,
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

func processPlayers(clubUrl string) (players []gin.H, _ error) {
	res, err := http.Get(clubUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
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
		players = append(players, gin.H{
			"id":  id,
			"url": BaseUrl + a.AttrOr("href", ""),
		})
	})
	return players, nil
}
