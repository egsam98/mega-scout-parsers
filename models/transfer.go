package models

import "time"

type Transfer struct {
	TransferType uint
	Date         time.Time
	FromTeam     string
	ToTeam       string
	Cost         string
	Fee          string
}
