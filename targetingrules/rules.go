package targetingrules

import (
	"fmt"

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

func New(rule string, value string) (TargetingRule, error) {
	if rule == "age" {
		return NewAgeTargetingRule(value), nil
	}

	return nil, fmt.Errorf("invalid rule: %s", rule)
}