package conditions

import (
	"testing"
)

func TestConditionBuilder_common(t *testing.T) {
	tests := []struct {
		name        string
		testBuilder *ConditionBuilder
		expected    string
	}{
		{
			name:        "simple",
			testBuilder: Par(Cond("a", "=", "b")),
			expected:    "(a = b)",
		},
		{
			name:        "harder",
			testBuilder: Par(Eq("a", "b").And().Ge("c", "10")).Or().Lt("a", "1"),
			expected:    "(a = b AND c >= 10) OR a < 1",
		},
		{
			name:        "is not null",
			testBuilder: Eq("a", "b").And().IsNotNull("c"),
			expected:    "a = b AND c IS NOT NULL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.testBuilder.String(); got != tt.expected {
				t.Errorf("And() = \"%v\", expected \"%v\"", got, tt.expected)
			}
		})
	}
}
