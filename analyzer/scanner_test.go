package analyzer

import (
	"testing"
)

func TestGetGrade(t *testing.T) {
	tests := []struct {
		tokens   int
		expected string
	}{
		{0, "A"},
		{1999, "A"},
		{2000, "B"},
		{9999, "B"},
		{10000, "C"},
		{29999, "C"},
		{30000, "D"},
		{99999, "D"},
		{100000, "F"},
		{500000, "F"},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.tokens)), func(t *testing.T) {
			result := getGrade(tt.tokens)
			if result != tt.expected {
				t.Errorf("getGrade(%d) = %s; want %s", tt.tokens, result, tt.expected)
			}
		})
	}
}

func TestGetGradeIndex(t *testing.T) {
	tests := []struct {
		grade    string
		expected int
	}{
		{"A", 0},
		{"B", 1},
		{"C", 2},
		{"D", 3},
		{"F", 4},
		{"Z", -1}, // Invalid grade
	}

	for _, tt := range tests {
		t.Run(tt.grade, func(t *testing.T) {
			result := GetGradeIndex(tt.grade)
			if result != tt.expected {
				t.Errorf("GetGradeIndex(%s) = %d; want %d", tt.grade, result, tt.expected)
			}
		})
	}
}
