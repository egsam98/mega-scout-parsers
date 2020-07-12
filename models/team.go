package models

type Team struct {
	Id      int      `json:"id" example:"2410"`
	Url     string   `json:"url" example:"https://transfermarkt.com/zska-moskau/startseite/verein/2410/saison_id/2019"`
	Title   string   `json:"title" example:"CSKA Moscow"`
	Players []Player `json:"players"`
}
