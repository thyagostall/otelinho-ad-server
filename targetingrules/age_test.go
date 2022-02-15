package targetingrules

import (
	"testing"

	"thyago.com/otelinho/openrtb"
)

func TestAge(t *testing.T) {
	cases := []struct {
		input    string
		operator TargetingOperator
		age      uint
	}{
		{input: ">99", operator: GreaterThan, age: 99},
		{input: "<99", operator: LessThan, age: 99},
		{input: ">=99", operator: GreaterThanOrEqual, age: 99},
		{input: "<=99", operator: LessThanOrEqual, age: 99},
		{input: "==99", operator: Equal, age: 99},
		{input: "!=99", operator: NotEqual, age: 99},
	}

	for _, tt := range cases {
		rule := NewAgeTargetingRule(tt.input)

		if rule.Operator != tt.operator {
			t.Fatalf("Invalid operator")
		}

		if rule.Value != tt.age {
			t.Fatalf("Invalid value")
		}
	}
}

func TestShouldInclude(t *testing.T) {
	cases := []struct {
		rule           string
		yob            uint
		expectedResult bool
	}{
		{rule: ">10", yob: 1990, expectedResult: true},
		{rule: "<10", yob: 1990, expectedResult: false},
		{rule: "==10", yob: 1990, expectedResult: false},
		{rule: "<=10", yob: 1990, expectedResult: false},
		{rule: ">=10", yob: 1990, expectedResult: true},
	}

	for _, tt := range cases {
		rule := NewAgeTargetingRule(tt.rule)

		result := rule.ShouldInclude(createBidRequest(tt.yob))
		if result != tt.expectedResult {
			t.Fatalf("Unexpected result: %v, should be %v", result, tt.expectedResult)
		}
	}
}

func createBidRequest(yob uint) *openrtb.BidRequest {
	return &openrtb.BidRequest{
		User: openrtb.User{YOB: yob},
	}
}
