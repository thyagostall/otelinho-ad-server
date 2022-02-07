package campaign

import "time"

type Campaign struct {
	ID        int
	Creative  string
	StartDate time.Time
	EndDate   time.Time
	Goal      uint
	MaxBid    float64
}
