package models

type Trophy struct {
	Title        string  `json:"title" example:"Russian champion"`
	SeasonPeriod string  `json:"season_period" example:"2020"`
	Team         *int    `json:"team" example:"964"`        // nullable
	Undetected   *string `json:"undetected" example:"null"` //nullable, не распознанный текст на событие/команду
	Event        *string `json:"event" example:"null"`      // nullable
	EventUrl     *string `json:"event_url" example:"null"`  // nullable
}
