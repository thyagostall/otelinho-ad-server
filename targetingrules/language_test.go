package targetingrules

import (
	"testing"

	"thyago.com/otelinho/openrtb"
)

func TestNewTargeting(t *testing.T) {
	cases := []struct {
		value            string
		expectedOperator TargetingOperator
		expectedValues   []string
	}{
		{"in[en,pt,es]", In, []string{"en", "pt", "es"}},
		{"in[en, pt, es]", In, []string{"en", "pt", "es"}},
		{"in[    en  ,    pt   , es    ]", In, []string{"en", "pt", "es"}},
		{"!in[en,pt,es]", In, []string{"en", "pt", "es"}},
		{"!in[en, pt, es]", In, []string{"en", "pt", "es"}},
		{"!in[    en  ,    pt   , es    ]", In, []string{"en", "pt", "es"}},
		{"in[en]", In, []string{"en"}},
		{"!in[en]", In, []string{"en"}},
	}

	for _, tt := range cases {
		rule := NewLanguageTargetingRule(tt.value)
		if !sliceEq(rule.Values, tt.expectedValues) {
			t.Fatalf("Expected %v, got %v (rule: %s)", tt.expectedValues, rule.Values, tt.value)
		}
	}
}

func TestLanguageShouldInclude(t *testing.T) {
	cases := []struct {
		rule           string
		value          string
		expectedResult bool
	}{
		{"in[en,pt,es]", "en", true},
		{"in[en,pt,es]", "cz", false},
		{"!in[en,pt,es]", "en", false},
		{"!in[en,pt,es]", "cz", true},
	}

	for _, tt := range cases {
		rule := NewLanguageTargetingRule(tt.rule)
		actualResult := rule.ShouldInclude(createLanguageBidRequest(tt.value))
		if actualResult != tt.expectedResult {
			t.Fatalf("Expected %v, got %v (rule: %s)", tt.expectedResult, actualResult, tt.rule)
		}
	}
}

func createLanguageBidRequest(language string) *openrtb.BidRequest {
	return &openrtb.BidRequest{
		Device: openrtb.Device{Language: language},
	}
}

func sliceEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
