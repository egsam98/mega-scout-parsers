package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func Trophies(personUrl string) (trophies []models.Trophy, _ error) {
	url := strings.ReplaceAll(personUrl, "profil", "erfolge")
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("div.box").Slice(0, -1).Each(func(_ int, box *goquery.Selection) {
		title := strings.Trim(box.Find("div[class$='header']").Text(), "\n\t ")
		title = regexp.MustCompile(`\d+x `).ReplaceAllString(title, "")
		box.Find("tr").Each(func(_ int, tr *goquery.Selection) {
			tds := tr.Find("td")
			a := tds.Find("a").Last()

			var team *int
			var undetected *string
			var event *string
			var eventUrl *string
			if a.Size() != 0 {
				text := a.Text()
				result, err := strconv.Atoi(a.AttrOr("id", ""))
				if err == nil {
					team = &result
				} else {
					result, exists := a.Attr("href")
					if exists && result != "" {
						eventUrl = new(string)
						*eventUrl = BaseUrl + result
						event = new(string)
						*event = text
					} else {
						undetected = &text
					}
				}
			}

			trophies = append(trophies, models.Trophy{
				Title:        title,
				SeasonPeriod: tds.First().Text(),
				Team:         team,
				Undetected:   undetected,
				Event:        event,
				EventUrl:     eventUrl,
			})
		})
	})
	return trophies, nil
}
