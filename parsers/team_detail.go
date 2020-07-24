package parsers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/message"
	"strconv"
	"strings"
	"time"
)

func TeamDetail(teamUrl string) (*models.TeamDetail, error) {
	foundedFuture := founded(teamUrl)

	doc, fetchHtmlErr := utils.FetchHtml(teamUrl)
	if fetchHtmlErr != nil {
		return nil, fetchHtmlErr
	}

	countryIdStr, exists := doc.Find("#land_select_breadcrumb > option").First().Attr("value")
	if !exists {
		return nil, fmt.Errorf("%s: country is absent", teamUrl)
	}

	var logo *string
	if result, exists := doc.Find(".dataBild > img").First().Attr("src"); exists {
		logo = &result
	}

	msg := <-foundedFuture
	if msg.IsError() {
		return nil, msg.Error
	}

	var founded *string
	if result, ok := msg.Data.(string); ok {
		founded = &result
	}

	country, err := strconv.Atoi(countryIdStr)
	if err != nil {
		return nil, err
	}
	return &models.TeamDetail{
		Country: country,
		Logo:    logo,
		Founded: founded,
	}, nil
}

func founded(teamUrl string) chan message.Message {
	future := make(chan message.Message, 1)
	url := strings.ReplaceAll(teamUrl, "startseite", "datenfakten")
	go func() {
		doc, fetchHtmlErr := utils.FetchHtml(url)
		if fetchHtmlErr != nil {
			future <- message.Error(fetchHtmlErr)
			return
		}

		foundedStr := doc.Find("table.profilheader tr").FilterFunction(func(_ int, tr *goquery.Selection) bool {
			return strings.Contains(strings.Trim(tr.Find("th").Text(), "\n\t "), "Founded")
		}).First().Find("td").Last().Text()

		if foundedStr == "" {
			future <- message.Nil()
			return
		}
		date, err := time.Parse("Jan 2, 2006", foundedStr)
		if err != nil {
			future <- message.Error(fmt.Errorf("%s: %v", teamUrl, err))
			return
		}
		future <- message.Ok(date.Format("02-01-2006"))
	}()
	return future
}
