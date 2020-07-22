package models

type Match struct {
	Id             int      `json:"id" example:"3381127"`
	Url            string   `json:"url" example:"https://www.transfermarkt.com/smena-saturn-st-petersburg_zenit-st-petersburg/index/spielbericht/3381127"`
	CompetitionUrl string   `json:"competition_url" example:"https://transfermarkt.com/russian-cup/startseite/pokalwettbewerb/RUP/saison_id/1993"`
	Round          string   `json:"round" example:"Third Round"`
	EventDatetime  string   `json:"event_datetime" example:"28-05-1993 12:00"`
	HomeTeam       int      `json:"home_team" example:"15170"`
	AwayTeam       int      `json:"away_team" example:"964"`
	HomeTeamScore  *int     `json:"home_team_score" example:"5"` // nullable
	AwayTeamScore  *int     `json:"away_team_score" example:"2"` // nullable
	Venue          string   `json:"venue" example:"Akademia Zenit"`
	LineUps        []LineUp `json:"line_ups"`
}

type LineUp struct {
	Team      int     `json:"team" example:"15170"`
	Formation *string `json:"formation" example:"4-5-1"` // nullable
	CoachId   int     `json:"coach_id" example:"77618"`
	CoachUrl  string  `json:"coach_url" example:"https://transfermarkt.com/viktor-vinogradov/profil/trainer/77618"`
	Players   []PlayerLineUp
}

type PlayerLineUp struct {
	Id     int  `json:"id" example:"751484"`
	Number *int `json:"number" example:"1"`           // nullable
	Type   int  `json:"type" enums:"0,1" example:"0"` // 0 - основной состав, 1 - в запасе
}

type PlayerStats struct {
	MatchId           int  `json:"match_id"`
	SeasonPeriod      int  `json:"season_period"`
	Goals             int  `json:"goals"`
	Assists           int  `json:"assists"`
	OwnGoals          int  `json:"own_goals"`
	YellowCards       int  `json:"yellow_cards"`
	SecondYellowCards int  `json:"second_yellow_cards"`
	RedCards          int  `json:"red_cards"`
	SubstitutionOn    *int `json:"substitution_on"`
	SubstitutionOff   *int `json:"substitution_off"`
	MinutesPlayed     int  `json:"minutes_played"`
}
