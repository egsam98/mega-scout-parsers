package parsers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/errors"
	"github.com/egsam98/MegaScout/utils/message"
	"github.com/egsam98/MegaScout/utils/slices"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var seasons = Seasons()

func Matches(teamUrl string) (matches []models.Match, _ error) {
	matchesChan := make(chan message.Message)
	for _, season := range seasons {
		go processSeasons(teamUrl, season, matchesChan)
	}

	for i := 0; i < len(seasons); i++ {
		msg := <-matchesChan
		if msg.IsError() {
			return nil, msg.Error
		}
		//fmt.Println(len(msg.Data.([]models.Match)))
		matches = append(matches, msg.Data.([]models.Match)...)
	}
	return matches, nil
}

func processSeasons(teamUrl string, season models.Season, matchesChan chan message.Message) {
	matchUrl := strings.ReplaceAll(teamUrl, "startseite", "spielplan") + "/saison_id/" +
		strconv.Itoa(season.Period)
	doc, err := utils.FetchHtml(matchUrl)
	if err != nil {
		matchesChan <- message.Error(err)
		return
	}

	matches := make([]models.Match, 0)
	var innerError error
	doc.Find("a.ergebnis-link").EachWithBreak(func(_ int, a *goquery.Selection) bool {
		if title, _ := a.Attr("title"); title != "Match report" {
			return true
		}
		href, exists := a.Attr("href")
		if !exists {
			innerError = err
			return false
		}
		matchUrl := "https://www.transfermarkt.com" + href
		match, err := matchInfo(matchUrl)
		if err != nil {
			innerError = err
			return false
		}
		matches = append(matches, *match)
		return true
	})

	if innerError != nil {
		matchesChan <- message.Error(innerError)
		return
	}
	matchesChan <- message.Ok(matches)
}

func matchInfo(matchUrl string) (*models.Match, error) {
	doc, fetchHtmlErr := utils.FetchHtml(matchUrl)
	if fetchHtmlErr != nil {
		return nil, fetchHtmlErr
	}

	text := doc.Find("div.ergebnis-wrap .sb-endstand").Text()

	scores := make([]*int, 2)
	for i, result := range strings.Split(regexp.MustCompile(`\d:\d`).FindString(text), ":") {
		if result != "" && result != "-" {
			score, err := strconv.Atoi(result)
			if err != nil {
				return nil, fmt.Errorf("%s: %v", matchUrl, err)
			}
			scores[i] = &score
		}
	}

	datum := doc.Find("p.sb-datum")
	formattedTime := regexp.MustCompile(`\d{1,2}:\d{1,2} [AaPp][Mm]`).FindString(datum.Text())
	formattedDatetime := strings.Trim(datum.Find("a").Eq(1).Text(), "\n\t ") + " " + formattedTime
	datetime, err := changeFormat(formattedDatetime)
	if err != nil {
		formattedDatetime := strings.Trim(datum.Find("a").First().Text(), "\n\t ") + " " + formattedTime
		datetime, err = changeFormat(formattedDatetime)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", matchUrl, err)
		}
	}

	teams := make([]int, 0)
	lineUps := make([]models.LineUp, 0)

	var innerErr error
	doc.Find(".box > .large-6.columns").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		team, err := strconv.Atoi(s.Find("nobr > a").First().AttrOr("id", ""))
		if err != nil {
			innerErr = fmt.Errorf("%s: %v", matchUrl, err)
			return false
		}
		teams = append(teams, team)
		coachUrl, _ := s.Find("a").Last().Attr("href")
		coachId, err := strconv.Atoi(slices.String_Last(strings.Split(coachUrl, "/")))
		if err != nil {
			innerErr = fmt.Errorf("%s: %v", matchUrl, err)
			return false
		}

		var formation *string
		if result := regexp.MustCompile(`\d-\d-\d`).FindString(s.Find("div.large-7").Text()); result != "" {
			formation = &result
		}

		lineUps = append(lineUps, models.LineUp{
			Team:      team,
			Formation: formation,
			CoachId:   coachId,
			CoachUrl:  BaseUrl + coachUrl,
		})
		return true
	})

	if innerErr != nil {
		return nil, innerErr
	}

	lineUpsUrl, exists := doc.Find("#line-ups > a").First().Attr("href")
	if exists && len(lineUps) != 0 {
		if err := processLineUps(matchUrl, BaseUrl+lineUpsUrl, &lineUps); err != nil {
			return nil, err
		}
	}

	id, err := strconv.Atoi(slices.String_Last(strings.Split(matchUrl, "/")))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", matchUrl, err)
	}

	competitionHref, exists := doc.Find(".spielername-profil a").First().Attr("href")
	if !exists {
		return nil, fmt.Errorf("%s: competition href is absent", matchUrl)
	}

	return &models.Match{
		Id:             id,
		Url:            matchUrl,
		CompetitionUrl: BaseUrl + competitionHref,
		Round:          strings.Trim(strings.Split(datum.Text(), "|")[0], "\n\t "),
		EventDatetime:  datetime,
		HomeTeam:       teams[0],
		AwayTeam:       teams[1],
		HomeTeamScore:  scores[0],
		AwayTeamScore:  scores[1],
		Venue:          doc.Find("span.hide-for-small a").Text(),
		LineUps:        lineUps,
	}, nil
}

func processLineUps(matchUrl, lineUpsUrl string, lineUps *[]models.LineUp) error {
	doc, err := utils.FetchHtml(lineUpsUrl)
	if err != nil {
		if _, ok := err.(*errors.TransfermarktError); ok {
			return nil
		}
		return err
	}

	var innerErr error
	doc.Find(".row.sb-formation").Slice(0, -1).EachWithBreak(func(typeLineUp int, e *goquery.Selection) bool {
		e.Find(".columns").EachWithBreak(func(i int, col *goquery.Selection) bool {
			playerLineUps := make([]models.PlayerLineUp, 0)
			col.Find("table.items > tbody > tr").EachWithBreak(func(_ int, tr *goquery.Selection) bool {
				tds := tr.Find("td")
				id, err := strconv.Atoi(tds.Eq(1).Find("a").First().AttrOr("id", ""))
				if err != nil {
					innerErr = fmt.Errorf("%s: %v", matchUrl, err)
					return false
				}

				var number *int
				numberStr := strings.Trim(tds.First().Text(), "\n\t ")
				if numberStr != "-" {
					result, err := strconv.Atoi(numberStr)
					if err != nil {
						innerErr = fmt.Errorf("%s: %v", matchUrl, err)
						return false
					}
					number = &result
				}

				playerLineUps = append(playerLineUps, models.PlayerLineUp{
					Id:     id,
					Number: number,
					Type:   typeLineUp,
				})
				return true
			})
			(*lineUps)[i].Players = append((*lineUps)[i].Players, playerLineUps...)

			if innerErr != nil {
				return false
			}
			return true
		})

		if innerErr != nil {
			return false
		}
		return true
	})

	return innerErr
}

func changeFormat(formattedDatetime string) (string, error) {
	result, err := time.Parse("Mon, 1/2/06 3:04 PM", formattedDatetime)
	if err != nil {
		return "", err
	}
	return result.Format("02-01-2006 15:04"), nil
}
