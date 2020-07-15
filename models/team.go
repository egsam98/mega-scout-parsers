package models

type Team struct {
	Id      int      `json:"id" example:"2410"`
	Url     string   `json:"url" example:"https://transfermarkt.com/zska-moskau/startseite/verein/2410/saison_id/2019"`
	Title   string   `json:"title" example:"CSKA Moscow"`
	Players []Player `json:"players"`
}

type TeamDetail struct {
	Country int     `json:"country" example:"141"`
	Logo    *string `json:"logo" example:"https://tmssl.akamaized.net/images/wappen/head/16704.png?lm=1499524238"` // nullable
	Founded *string `json:"founded" swaggertype:"string" example:"22-02-2008"`                                     // nullable
}
