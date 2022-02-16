package targetingrules

import (
	"testing"
	"time"

	"thyago.com/otelinho/openrtb"
)

func TestNewAgeRule(t *testing.T) {
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

func TestAgeShouldInclude(t *testing.T) {
	now := uint(time.Now().Year())

	cases := []struct {
		rule           string
		yob            uint
		expectedResult bool
	}{
		{rule: ">30", yob: now - 31, expectedResult: true},
		{rule: "<30", yob: now - 29, expectedResult: true},
		{rule: "==30", yob: now - 30, expectedResult: true},
		{rule: "!=30", yob: now, expectedResult: true},
		{rule: "<=30", yob: now - 30, expectedResult: true},
		{rule: ">=30", yob: now - 30, expectedResult: true},

		{rule: ">30", yob: now - 29, expectedResult: false},
		{rule: "<30", yob: now - 31, expectedResult: false},
		{rule: "==30", yob: now, expectedResult: false},
		{rule: "!=30", yob: now - 30, expectedResult: false},
		{rule: "<=30", yob: now - 31, expectedResult: false},
		{rule: ">=30", yob: now - 29, expectedResult: false},
	}

	for _, tt := range cases {
		rule := NewAgeTargetingRule(tt.rule)

		result := rule.ShouldInclude(createAgeBidRequest(tt.yob))
		if result != tt.expectedResult {
			t.Fatalf("Unexpected result: %v, should be %v (Rule: %s, YOB: %d)", result, tt.expectedResult, tt.rule, tt.yob)
		}
	}
}

func createAgeBidRequest(yob uint) *openrtb.BidRequest {
	return &openrtb.BidRequest{
		User: openrtb.User{YOB: yob},
	}
}
