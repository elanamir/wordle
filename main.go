package main

import (
	"fmt"
	str "strings"
)

const WORDLEN = 5

type Pattern []int

var AllPatterns []Pattern

func main() {
	// Allwords = read_words(fname)

	AllPatterns = make([]Pattern, 0)

	// build_strategy(Allwords)
	BuildPattern(0, Pattern{})
	// fmt.Println(patterns)

}

type StrategyElement struct {
	Pattern   Pattern
	NextStage *StrategyStage
}

type StrategyStage struct {
	Dictionary []string
	Guess      string
	Patterns   []StrategyElement
}

func InitStrategy(allwords []string) StrategyStage {
	root := StrategyStage{Dictionary: allwords}
	BuildStrategy(&root, allwords)
	PrintStrategy(&root, 0, nil)
	return root
}

func PrintStrategy(s *StrategyStage, depth int, pattern Pattern) {
	if s != nil {
		for i := 0; i < depth; i += 1 {
			fmt.Print("  ")
		}
		if pattern != nil {
			fmt.Print(pattern)
		}
		if len(s.Dictionary) != 0 {
			fmt.Print(" Depth:", depth, " Guess: ", s.Guess, " DictLen: ", len(s.Dictionary))

			if len(s.Dictionary) == 1 {
				fmt.Print(" Final Word: ", s.Dictionary[0])
			}

			fmt.Println("")

			if len(s.Patterns) != 0 {
				for _, p := range s.Patterns {
					if p.NextStage != nil && len(p.NextStage.Dictionary) != 0 {
						PrintStrategy(p.NextStage, depth+1, p.Pattern)
					}
				}
			}
		}
	}
}

func BuildPattern(position int, pattern Pattern) {
	if position == WORDLEN {
		AllPatterns = append(AllPatterns, pattern)
		return
	}
	for _, i := range []int{0, 1, 2} {
		p := append(pattern, i)
		BuildPattern(position+1, p)
	}
}

func BuildStrategy(s *StrategyStage, allwords []string) {
	minmax_len := 1000000 //XXX

	// fmt.Println("Entered")
	if len(s.Dictionary) < 2 {
		// fmt.Println("returning")
		return
	}

	// fmt.Println("strategy stage", s)

	var best StrategyStage
	for _, guess := range allwords {
		// fmt.Println("guess", guess)
		cur := StrategyStage{Guess: guess}
		max_pattern_len := 0

		// Determine how many words from the dictionary match for each pattern
		for _, pattern := range AllPatterns {
			// fmt.Println("pattern", pattern)
			se := StrategyElement{Pattern: pattern}
			se.NextStage = &StrategyStage{}

			for _, word := range s.Dictionary {
				// fmt.Println("word", word)
				if PatternMatch(guess, word, pattern) {
					fmt.Println("match", guess, word, pattern)

					// Add the word to the pattern's dictionary
					se.NextStage.Dictionary = append(se.NextStage.Dictionary, word)
					fmt.Println("Se1", se, se.NextStage)
				}
			}
			fmt.Println("guess", guess, max_pattern_len)

			if max_pattern_len < len(se.NextStage.Dictionary) {
				max_pattern_len = len(se.NextStage.Dictionary)

			}
			fmt.Println("Se", se, se.NextStage)
			cur.Patterns = append(cur.Patterns, se)
		}

		// If the max is less than the current min, this is the best strategy
		if minmax_len > max_pattern_len {
			minmax_len = max_pattern_len
			best = cur
			fmt.Println("best guess", best.Guess, max_pattern_len)
		}
	}

	s.Guess = best.Guess
	s.Patterns = best.Patterns

	for _, p := range s.Patterns {
		BuildStrategy(p.NextStage, allwords)
	}
}

// XXX May be more efficient ways to do this, but let's just get started
// Does word match guess according to pattern?
func PatternMatch(guess string, word string, pattern []int) bool {
	for ip, cp := range pattern {
		cg := guess[ip]
		// fmt.Println("ip", ip, "cp", cp, "cg", string(cg), "word", word, str.IndexByte(word, cg))
		idx := str.IndexByte(word, cg)
		switch cp {
		case 0:
			if idx >= 0 {
				return false
			}
		case 1:
			if idx < 0 || idx == ip { // -1 or positive and equal
				return false
			}
		case 2:
			if idx < 0 || idx != ip {
				return false
			}
		default:
			// error
		}
	}
	return true
}

/**
 * XXX this isn't right yet...
 * For each stage, we pick the guess from all guesses that gave the minimum number of maximum word sets across all responses, and that's ths strategy.
 *
S = ALLWORDS
lastmin = infinity
while length S > 1 do: # XXX is this condition?
	foreach guess in S do:
		lastmax = 0
		foreach response in (00000..22222) do:
			A[response] = set of words v from S where v satisfies guess and response
			m = length of A[response]
			if (m > lastmax) lastmax = m
		end

		if (lastmax < lastmin) { lastmin = lastmax, bestguess = guess, S = A }
	end
	# need to structure strategy
end

**/
