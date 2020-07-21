package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/slices"
	strings2 "github.com/egsam98/MegaScout/utils/strings"
	"regexp"
	"strings"
)

func MatchEvents(matchUrl string) (matchEvents []interface{}, _ error) {
	doc, err := utils.FetchHtml(matchUrl)
	if err != nil {
		return nil, err
	}

	matchEventChan := make(chan interface{})

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
		matchEvents = append(matchEvents, <-matchEventChan)
	}
	return matchEvents, nil
}

func processGoals(lis *goquery.Selection, matchEventChan chan interface{}) {
	lis.Each(func(_ int, li *goquery.Selection) {
		action := li.Find(".sb-aktion-aktion")
		goalAndAssist := [2]*string{}
		for i, splitted := range strings.Split(action.Text(), "Assist:") {
			if result := strings.Trim(strings.Split(splitted, ",")[1], "\n\t "); result != "" {
				goalAndAssist[i] = &result
			}
		}

		goalAndAssistPlayer := [2]*int{}
		action.Find("a").Each(func(i int, a *goquery.Selection) {
			id := strings2.ToInt(a.AttrOr("id", ""), false)
			goalAndAssistPlayer[i] = &id
		})

		matchEventChan <- models.NewGoal(
			li,
			*goalAndAssistPlayer[0],
			*goalAndAssist[0],
			goalAndAssistPlayer[1],
			goalAndAssist[1],
		)
	})
}

func processSubstitutions(lis *goquery.Selection, matchEventChan chan interface{}) {
	lis.Each(func(_ int, li *goquery.Selection) {
		ids := utils.NewSet()
		li.Find(".sb-aktion-aktion a").Each(func(_ int, a *goquery.Selection) {
			id := strings2.ToInt(a.AttrOr("id", ""), false)
			ids.Add(id)
		})

		slice := ids.Slice()
		matchEventChan <- models.NewSubstitution(
			li,
			slice[0].(int),
			slice[1].(int),
		)
	})
}

func processCards(lis *goquery.Selection, matchEventChan chan interface{}) {
	lis.Each(func(_ int, li *goquery.Selection) {
		action := li.Find(".sb-aktion-aktion")
		player := strings2.ToInt(action.Find("a").First().AttrOr("id", ""), false)
		info := slices.String_Last(strings.Split(action.Text(), "\t"))
		info = strings.Trim(regexp.MustCompile(`\d.`).ReplaceAllString(info, ""), "\n\t ")

		matchEventChan <- models.NewCard(
			li,
			player,
			info,
		)
	})
}

func processPenalty(lis *goquery.Selection, matchEventChan chan interface{}) {
	lis.Each(func(_ int, li *goquery.Selection) {
		player := strings2.ToInt(li.Find(".sb-aktion-aktion > a").First().AttrOr("id", ""), false)
		matchEventChan <- models.NewPenalty(
			li,
			player,
			li.Find(".sb-11m-tor").Length() > 0,
		)
	})
}

func findTable(tables *goquery.Selection, header string) *goquery.Selection {
	return tables.FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.Contains(s.Find(".header").Text(), header)
	}).First().Find("li")
}
