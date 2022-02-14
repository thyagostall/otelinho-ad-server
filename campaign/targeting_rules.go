package campaign

import (
	"time"

	"thyago.com/otelinho/openrtb"
)

type TargetingOperator uint

const (
	Equal TargetingOperator = iota
	NotEqual
	LessThan
	LessThanOrEqual
	GreaterThan
	GreaterThanOrEqual
)

type TargetingRule interface {
	ShouldInclude(*openrtb.BidRequest) bool
}

type AgeTargetingRule struct {
	Operator TargetingOperator
	Value    uint
}

func (r AgeTargetingRule) ShouldInclude(candidate *openrtb.BidRequest) bool {
	year := time.Now().UTC().Year()
	age := year - candidate.User.YOB

	if r.Operator == Equal && candidate.User.YOB == age {
		return true
	} else if r.Operator == NotEqual && candidate.User.YOB != age {
		return true
	} else if r.Operator == LessThan && candidate.User.YOB < age {
		return true
	} else if r.Operator == LessThanOrEqual && candidate.User.YOB <= age {
		return true
	} else if r.Operator == GreaterThan && candidate.User.YOB > age {
		return true
	} else if r.Operator == GreaterThanOrEqual && candidate.User.YOB >= age {
		return true
	} else {
		return false
	}
}
