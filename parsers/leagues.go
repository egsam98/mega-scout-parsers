package parsers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/slices"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

const BaseUrl = "https://transfermarkt.com"

func Leagues(countryId, seasonPeriod int) ([]models.League, error) {
	url := fmt.Sprintf("%s/wettbewerbe/national/wettbewerbe/%d?saison_id=%d", BaseUrl, countryId, seasonPeriod)
	doc, err := utils.RetryFetchHtml(url, 5)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tier := ""
	leagues := make([]models.League, 0)
	doc.Find("#yw1 tbody > tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		if strings.Contains(tr.Text(), "Cup") {
			return false
		}
		td := tr.Find("td").First()
		if strings.Contains(td.AttrOr("class", ""), "extrarow") {
			tier = td.Text()
			return true
		}
		if _, exists := tr.Attr("class"); !exists {
			return true
		}
		a := td.Find("a").Last()
		href, _ := a.Attr("href")

		regex, _ := regexp.Compile(`/saison_id/\d+/?`)
		urlWithoutSeasonId := regex.Split(href, 2)[0]
		idStr := slices.String_Last(strings.Split(urlWithoutSeasonId, "/"))
		leagues = append(leagues, models.League{
			Id:       generateId(idStr),
			Url:      BaseUrl + urlWithoutSeasonId,
			Position: tier,
		})
		return true
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
