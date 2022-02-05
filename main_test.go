package main

import (
	"fmt"
	"testing"
)

func TestPatternMatch(t *testing.T) {
	var tests = []struct {
		guess, word string
		pattern     Pattern
		want        bool
	}{
		{"abate", "aahed", []int{2, 0, 1, 0, 0}, false},
		{"abate", "aahed", []int{2, 0, 1, 0, 1}, true},
		{"scare", "stake", []int{2, 0, 2, 0, 2}, true},
		{"blink", "plink", []int{0, 0, 0, 0, 0}, false},
		{"scare", "stake", []int{2, 1, 2, 1, 1}, false},
		{"scare", "crane", []int{0, 1, 2, 1, 2}, true},
		{"scare", "cribs", []int{1, 1, 0, 1, 0}, true},
		{"plink", "crave", []int{0, 0, 0, 1, 0}, false},
		{"aback", "brake", []int{1, 1, 2, 0, 1}, false},
		{"aback", "brake", []int{0, 1, 2, 0, 1}, true},
		{"abaca", "pleat", []int{1, 0, 0, 0, 0}, true},
		{"abaca", "pleat", []int{1, 0, 1, 0, 1}, false},
		{"aback", "abaca", []int{1, 1, 1, 1, 0}, false},
		{"scare", "crack", []int{0, 1, 2, 1, 0}, true}, // XXX what is right pattern?
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%+v", tt)
		t.Run(testname, func(t *testing.T) {
			ans := PatternMatch(tt.guess, tt.word, tt.pattern)
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

func TestBuildStrategy(t *testing.T) {
	var allwords = []string{
		"scare",
		"blink",
		"crave",
		"brave",
		"crass",
		"bored",
		"bland",
		"stand",
		"stink",
		"drink",
		"plink",
		"brass",
		"dress",
		"brand",
	}
	var tests = []struct {
		allwords []string
		want     bool
	}{
		{allwords, true},
	}
	AllPatterns = make([]Pattern, 0)
	BuildPattern(0, Pattern{}, 5) // XXX

	for _, tt := range tests {
		testname := fmt.Sprintf("%+v", tt)
		t.Run(testname, func(t *testing.T) {
			InitStrategy(allwords)
		})
	}
}
