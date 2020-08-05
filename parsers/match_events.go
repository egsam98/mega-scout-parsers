package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/message"
	"github.com/egsam98/MegaScout/utils/slices"
	. "github.com/go-errors/errors"
	"regexp"
	"strconv"
	"strings"
)

func MatchEvents(matchUrl string) (matchEvents []interface{}, _ *Error) {
	doc, err := utils.FetchHtml(matchUrl)
	if err != nil {
		return nil, New(err)
	}

	matchEventChan := make(chan message.Message)

	tables := doc.Find("[class='row']")
	goals := findTable(tables, "Goals")
	substitutions := findTable(tables, "Substitutions")
	cards := findTable(tables, "Cards")
	penalty := findTable(tables, "Penalty shoot-out")

	go processGoals(goals, matchEventChan)
	go processSubstitutions(substitutions, matchEventChan)
	go processCards(cards, matchEventChan)
	go processPenalty(penalty, matchEventChan)

	for i := 0; i < goals.Length()+substitutions.Length()+cards.Length()+penalty.Length(); i++ {
		msg := <-matchEventChan
		if msg.IsError() {
			return nil, msg.Error
		}
		matchEvents = append(matchEvents, msg.Data)
	}
	return matchEvents, nil
}

func processGoals(lis *goquery.Selection, matchEventChan chan message.Message) {
	lis.Each(func(_ int, li *goquery.Selection) {
		action := li.Find(".sb-aktion-aktion")
		goalAndAssist := [2]*string{}
		for i, splitted := range strings.Split(action.Text(), "Assist:") {
			if result := strings.Trim(strings.Split(splitted, ",")[1], "\n\t "); result != "" {
				goalAndAssist[i] = &result
			}
		}

		goalAndAssistPlayer := [2]*int{}
		var innerErr *Error
		action.Find("a").EachWithBreak(func(i int, a *goquery.Selection) bool {
			id, err := strconv.Atoi(a.AttrOr("id", ""))
			if err != nil {
				innerErr = New(err)
				return false
			}
			goalAndAssistPlayer[i] = &id
			return true
		})

		if innerErr != nil {
			matchEventChan <- message.Error(innerErr)
			return
		}

		goal, err := models.NewGoal(
			li,
			*goalAndAssistPlayer[0],
			*goalAndAssist[0],
			goalAndAssistPlayer[1],
			goalAndAssist[1],
		)
		if err != nil {
			matchEventChan <- message.Error(New(err))
			return
		}
		matchEventChan <- message.Ok(goal)
	})
}

func processSubstitutions(lis *goquery.Selection, matchEventChan chan message.Message) {
	lis.Each(func(_ int, li *goquery.Selection) {
		ids := utils.NewSet()
		var innerErr *Error
		li.Find(".sb-aktion-aktion a").EachWithBreak(func(_ int, a *goquery.Selection) bool {
			id, err := strconv.Atoi(a.AttrOr("id", ""))
			if err != nil {
				innerErr = New(err)
				return false
			}
			ids.Add(id)
			return true
		})

		if innerErr != nil {
			matchEventChan <- message.Error(innerErr)
			return
		}

		slice := ids.Slice()
		sub, err := models.NewSubstitution(
			li,
			slice[0].(int),
			slice[1].(int),
		)
		if err != nil {
			matchEventChan <- message.Error(New(err))
			return
		}
		matchEventChan <- message.Ok(sub)
	})
}

func processCards(lis *goquery.Selection, matchEventChan chan message.Message) {
	lis.Each(func(_ int, li *goquery.Selection) {
		action := li.Find(".sb-aktion-aktion")
		player, err := strconv.Atoi(action.Find("a").First().AttrOr("id", ""))
		if err != nil {
			matchEventChan <- message.Error(New(err))
			return
		}

		info := slices.String_Last(strings.Split(action.Text(), "\t"))
		info = strings.Trim(regexp.MustCompile(`\d.`).ReplaceAllString(info, ""), "\n\t ")
		card, err := models.NewCard(li, player, info)
		if err != nil {
			matchEventChan <- message.Error(New(err))
			return
		}
		matchEventChan <- message.Ok(card)
	})
}

func processPenalty(lis *goquery.Selection, matchEventChan chan message.Message) {
	lis.Each(func(_ int, li *goquery.Selection) {
		player, err := strconv.Atoi(li.Find(".sb-aktion-aktion > a").First().AttrOr("id", ""))
		if err != nil {
			matchEventChan <- message.Error(New(err))
			return
		}
		penalty, err := models.NewPenalty(
			li,
			player,
			li.Find(".sb-11m-tor").Length() > 0,
		)
		if err != nil {
			matchEventChan <- message.Error(New(err))
			return
		}
		matchEventChan <- message.Ok(penalty)
	})
}

func findTable(tables *goquery.Selection, header string) *goquery.Selection {
	return tables.FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.Contains(s.Find(".header").Text(), header)
	}).First().Find("li")
}
