package targetingrules

import (
	"regexp"
	"strings"

	"thyago.com/otelinho/openrtb"
)

type LanguageTargetingRule struct {
	Operator TargetingOperator
	Values   []string
}

var regexpLanguageOperatorValue *regexp.Regexp

func init() {
	regexpLanguageOperatorValue, _ = regexp.Compile(`(in|!in)\[((\S+|,|\s+)+)\]`)
}

func NewLanguageTargetingRule(rawValue string) LanguageTargetingRule {
	elements := regexpLanguageOperatorValue.FindStringSubmatch(rawValue)
	rawOperator := elements[1]
	rawValues := elements[2]

	var operator TargetingOperator
	if rawOperator == "in" {
		operator = In
	} else if rawOperator == "!in" {
		operator = NotIn
	}

	var values []string
	for _, elem := range strings.Split(rawValues, ",") {
		values = append(values, strings.TrimSpace(elem))
	}

	return LanguageTargetingRule{Operator: operator, Values: values}
}

func (r LanguageTargetingRule) ShouldInclude(candidate *openrtb.BidRequest) bool {
	if r.Operator == In {
		return inArray(candidate.Device.Language, r.Values)
	} else if r.Operator == NotIn {
		return !inArray(candidate.Device.Language, r.Values)
	}

	return false
}

func inArray(value string, values []string) bool {
	for _, v := range values {
		if value == v {
			return true
		}
	}
	return false
}
