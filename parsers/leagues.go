package parsers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/slices"
	. "github.com/go-errors/errors"
	"regexp"
	"strconv"
	"strings"
)

const BaseUrl = "https://transfermarkt.com"

func Leagues(countryId, seasonPeriod int) (leagues []models.League, _ *Error) {
	url := fmt.Sprintf("%s/wettbewerbe/national/wettbewerbe/%d?saison_id=%d", BaseUrl, countryId, seasonPeriod)
	doc, err := utils.FetchHtml(url)
	if err != nil {
		return nil, New(err)
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
		}

		regex, _ := regexp.Compile(`/saison_id/\d+/?`)
		urlWithoutSeasonId := regex.Split(href, 2)[0]
		idStr := slices.String_Last(strings.Split(urlWithoutSeasonId, "/"))
		leagues = append(leagues, models.League{
			Id:       generateId(idStr),
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
