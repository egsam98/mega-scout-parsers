package models

type Season struct {
	Period      int    `json:"period" example:"2019"`
	SeasonStart string `json:"season_start" example:"01-07-2019"`
	SeasonEnd   string `json:"season_end" example:"30-06-2020"`
}
