package targetingrules

import (
	"regexp"
	"strconv"
	"time"

	"thyago.com/otelinho/openrtb"
)

type AgeTargetingRule struct {
	Operator TargetingOperator
	Value    uint
}

var regexpOperatorValue *regexp.Regexp

func init() {
	regexpOperatorValue, _ = regexp.Compile(`(==|!=|<|<=|>|>=)(\d+)`)
}

func NewAgeTargetingRule(rawValue string) AgeTargetingRule {
	elements := regexpOperatorValue.FindStringSubmatch(rawValue)

	var operator TargetingOperator
	if elements[1] == "==" {
		operator = Equal
	} else if elements[1] == "!=" {
		operator = NotEqual
	} else if elements[1] == "<" {
		operator = LessThan
	} else if elements[1] == "<=" {
		operator = LessThanOrEqual
	} else if elements[1] == ">" {
		operator = GreaterThan
	} else if elements[1] == ">=" {
		operator = GreaterThanOrEqual
	}

	value, _ := strconv.ParseUint(elements[2], 10, 32)
	return AgeTargetingRule{Operator: operator, Value: uint(value)}
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