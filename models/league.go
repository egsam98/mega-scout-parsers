package models

type League struct {
	Id       int     `json:"id" example:"17201"`
	Url      string  `json:"url" example:"https://transfermarkt.com/premier-liga/startseite/wettbewerb/RU1"`
	Title    string  `json:"title" example:"Premier Liga"`
	Logo     *string `json:"logo" example:"https://tmssl.akamaized.net/images/logo/normal/ru1.png?lm=1582769594" extensions:"x-nullable"`
	Position string  `json:"position" example:"First Tier"`
}
