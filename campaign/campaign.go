package campaign

import (
	"time"

	"thyago.com/otelinho/targetingrules"
)

type Campaign struct {
	ID              int
	Creative        string
	StartDate       time.Time
	EndDate         time.Time
	Goal            uint
	MaxBid          float64
	RemainingBudget float64
	Budget          float64
	PacingFactor    int64
	Targeting       []targetingrules.TargetingRule
}
