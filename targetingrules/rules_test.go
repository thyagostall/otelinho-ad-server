package targetingrules

import "testing"

func TestNewWithValidRule(t *testing.T) {
	rule, _ := New("age", ">20")
	_, ok := rule.(AgeTargetingRule)

	if !ok {
		t.Fatalf("Invalid rule: %T, expected: %T\n", rule, AgeTargetingRule{})
	}
}

func TestNewWithInvalidRule(t *testing.T) {
	rule, err := New("invalid", "any value")

	if err == nil {
		t.Fatalf("Expected error\n")
	}

	if rule != nil {
		t.Fatalf("Shouldn't return any value on error")
	}
}
