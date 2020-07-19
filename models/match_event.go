package models

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

type AllMatchEventFields struct {
	Goal
	Penalty
	Card
	Substitution
}

type MatchEvent struct {
	Type string `json:"type" enums:"goal,penalty,substitution,card" example:"goal"`
	Team int    `json:"team" example:"2741"`
}

type Goal struct {
	MatchEvent
	Minute       int     `json:"minute" example:"71"`                   // goal
	GoalPlayer   int     `json:"goal_player" example:"121434"`          // goal
	GoalInfo     string  `json:"goal_info" example:"Right-footed shot"` // goal
	AssistPlayer *int    `json:"assist_player" example:"242567"`        // goal, nullable
	AssistInfo   *string `json:"assist_info" example:"Pass"`            // goal, nullable
}

func NewGoal(li *goquery.Selection, goalPlayer int, goalInfo string, assistPlayer *int, assistInfo *string) Goal {
	return Goal{
		MatchEvent: MatchEvent{
			Type: "goal",
			Team: team(li),
		},
		Minute:       minute(li),
		GoalPlayer:   goalPlayer,
		GoalInfo:     goalInfo,
		AssistPlayer: assistPlayer,
		AssistInfo:   assistInfo,
	}
}

type Penalty struct {
	MatchEvent
	Player  int  `json:"player" example:"184071"` // penalty
	Success bool `json:"success" example:"true"`  // penalty
}

func NewPenalty(li *goquery.Selection, player int, success bool) Penalty {
	return Penalty{
		MatchEvent: MatchEvent{
			Type: "penalty",
			Team: team(li),
		},
		Player:  player,
		Success: success,
	}
}

type Substitution struct {
	MatchEvent
	Minute    int `json:"minute" example:"71"`         // substitution
	PlayerIn  int `json:"player_in" example:"125624"`  // substitution
	PlayerOut int `json:"player_out" example:"956421"` // substitution
}

func NewSubstitution(li *goquery.Selection, playerIn, playerOut int) Substitution {
	return Substitution{
		MatchEvent: MatchEvent{
			Type: "substitution",
			Team: team(li),
		},
		Minute:    minute(li),
		PlayerIn:  playerIn,
		PlayerOut: playerOut,
	}
}

type Card struct {
	MatchEvent
	Minute int    `json:"minute" example:"71"`              // card
	Player int    `json:"player" example:"123456"`          // card
	Info   string `json:"info" example:"Yellow card, Foul"` // card
}

func NewCard(li *goquery.Selection, player int, info string) Card {
	return Card{
		MatchEvent: MatchEvent{
			Type: "card",
			Team: team(li),
		},
		Minute: minute(li),
		Player: player,
		Info:   info,
	}
}

func minute(li *goquery.Selection) int {
	style := li.Find(".sb-sprite-uhr-klein").First().AttrOr("style", "")
	xY := [2]int{}
	for i, pxStr := range strings.Split(style, " ")[1:] {
		px, err := strconv.Atoi(strings.Trim(pxStr, "px;"))
		if err != nil {
			panic(err)
		}
		xY[i] = px / -36
	}
	return 10*xY[1] + xY[0] + 1
}

func team(node *goquery.Selection) int {
	result, err := strconv.Atoi(node.Find(".sb-aktion-wappen > a").First().AttrOr("id", ""))
	if err != nil {
		panic(err)
	}
	return result
}
