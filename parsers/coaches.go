package parsers

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/message"
	"github.com/egsam98/MegaScout/utils/slices"
	"strconv"
	"strings"
)

func Coaches(matchUrl string) ([]models.Coach, error) {
	doc, err := utils.FetchHtml(matchUrl)
	if err != nil {
		return nil, err
	}

	ch := make(chan message.Message)
	doc.Find(".ersatzbank").Each(func(_ int, e *goquery.Selection) {
		href, exists := e.Find("a").Last().Attr("href")
		if !exists {
			panic(errors.New("coach href doesn't exist"))
		}
		go processCoach(BaseUrl+href, ch)
	})

	coaches := make([]models.Coach, 2)
	for i := 0; i < 2; i++ {
		msg := <-ch
		if msg.IsError() {
			return nil, msg.Error
		}
		coaches[i] = msg.Data.(models.Coach)
	}

	return coaches, nil
}

func processCoach(url string, ch chan<- message.Message) {
	_, err := utils.FetchHtml(url)
	if err != nil {
		ch <- message.Error(err)
		return
	}
	id, err := strconv.Atoi(slices.String_Last(strings.Split(url, "/")))
	if err != nil {
		panic(err)
	}
	ch <- message.Ok(models.Coach{
		Id:  id,
		Url: url,
	})
}
