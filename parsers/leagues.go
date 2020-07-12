package parsers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const BaseUrl = "https://transfermarkt.com"

func Leagues(countryId, seasonPeriod int) (leagues []models.League, _ error) {
	url := fmt.Sprintf("%s/wettbewerbe/national/wettbewerbe/%d?saison_id=%d", BaseUrl, countryId, seasonPeriod)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", url, err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	tier := ""
	doc.Find("#yw1 tbody > tr").Each(func(i int, tr *goquery.Selection) {
		if strings.Contains(tr.Text(), "Cup") {
			return
		}
		td := tr.Find("td").First()
		tdClass, _ := td.Attr("class")
		if strings.Contains(tdClass, "extrarow") {
			tier = td.Text()
			return
		}
		if _, exists := tr.Attr("class"); !exists {
			return
		}
		a := td.Find("a").Last()
		href, _ := a.Attr("href")

		logoStr, exists := td.Find("img").First().Attr("src")
		var logo *string
		if exists {
			logoStr = strings.ReplaceAll(logoStr, "tiny", "normal")
			logo = &logoStr
		} else {
			logo = nil
		}

		regex, _ := regexp.Compile(`/saison_id/\d+/?`)
		urlWithoutSeasonId := regex.Split(href, 2)[0]
		splitted := strings.Split(urlWithoutSeasonId, "/")
		leagues = append(leagues, models.League{
			Id:       generateId(splitted[len(splitted)-1]),
			Url:      BaseUrl + urlWithoutSeasonId,
			Title:    strings.Trim(td.Text(), "\t\n "),
			Logo:     logo,
			Position: tier,
		})
	})
	return leagues, nil
}

func generateId(id string) int {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	result := ""
	for _, str := range id {
		str := string(str)
		_, err := strconv.Atoi(str)
		if err != nil {
			result += strconv.Itoa(strings.Index(alphabet, str))
		} else {
			result += str
		}
	}
	resultNum, _ := strconv.Atoi(result)
	return resultNum
}
