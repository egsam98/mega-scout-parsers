package models

type Transfer struct {
	TransferType int     `json:"transfer_type" enums:"0,1" example:"0"` // 1 - арена, 0 - все остальное
	Date         *string `json:"date" example:"2017-01-01T00:00:00Z"`
	Season       string  `json:"season"`
	FromTeam     *string `json:"from_team" example:"28095"` // nullable
	ToTeam       string  `json:"to_team" example:"964"`     // nullable
	Cost         *string `json:"cost" example:"€2.00m"`     // nullable
	Fee          *string `json:"fee" example:"€3.50m"`      // nullable
}
