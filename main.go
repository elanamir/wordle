package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	str "strings"
)

type Pattern []int

var AllPatterns []Pattern

func main() {

	fnamePtr := flag.String("wordfile", "./words/5-3.txt", "wordlist file")
	nwordsPtr := flag.Int("nwords", -1, "number of words to process (all if -1)")
	wordlenPtr := flag.Int("wordlen", 5, "Length of words in word list")

	flag.Parse()

	AllPatterns = make([]Pattern, 0)
	BuildPattern(0, Pattern{}, *wordlenPtr)

	allwords := BuildWordList(*fnamePtr, *nwordsPtr)
	fmt.Println("Num Words: ", len(allwords))
	InitStrategy(allwords)
}

func BuildWordList(fname string, nwords int) []string {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("failed to open")

	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for n := 0; ((nwords == -1) || n < nwords) && scanner.Scan(); n += 1 {
		text = append(text, scanner.Text())
	}
	file.Close()

	return text
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
	maxdepth := -1
	PrintStrategy(&root, 0, nil, &maxdepth)
	fmt.Println("Max Depth: ", maxdepth)
	return root
}

func PrintStrategy(s *StrategyStage, depth int, pattern Pattern, maxdepth *int) {
	if s == nil {
		return
	}

	if depth > *maxdepth {
		*maxdepth = depth
	}

	for i := 0; i < depth; i += 1 {
		fmt.Print("  ")
	}

	if pattern != nil {
		fmt.Printf("d%d%v", depth, pattern)
	}
	if len(s.Dictionary) != 0 {

		if len(s.Dictionary) == 1 {
			fmt.Print(" Final Word: ", s.Dictionary[0])
		} else {
			fmt.Print(" Guess: ", s.Guess, " DictLen: ", len(s.Dictionary))
		}

		fmt.Println("")

		if len(s.Patterns) != 0 {
			for _, p := range s.Patterns {
				if p.NextStage != nil && len(p.NextStage.Dictionary) != 0 {
					PrintStrategy(p.NextStage, depth+1, p.Pattern, maxdepth)
				}
			}
		}
	}
}

func BuildPattern(position int, pattern Pattern, wordlen int) {
	if position == wordlen {
		AllPatterns = append(AllPatterns, pattern)
		return
	}
	for _, i := range []int{0, 1, 2} {
		p := append(pattern, i)
		BuildPattern(position+1, p, wordlen)
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
		cur := StrategyStage{Guess: guess}
		max_pattern_len := 0

		// Determine how many words from the dictionary match for each pattern
		for _, pattern := range AllPatterns {
			se := StrategyElement{Pattern: pattern}
			se.NextStage = &StrategyStage{}

			for _, word := range s.Dictionary {
				if word != guess && PatternMatch(guess, word, pattern) {

					// Add the word to the pattern's dictionary
					se.NextStage.Dictionary = append(se.NextStage.Dictionary, word)
				}
			}

			if max_pattern_len < len(se.NextStage.Dictionary) {
				max_pattern_len = len(se.NextStage.Dictionary)

			}
			cur.Patterns = append(cur.Patterns, se)
		}

		// If the max is less than the current min, this is the best strategy
		if minmax_len > max_pattern_len {
			minmax_len = max_pattern_len
			best = cur
			// fmt.Println("best guess:", best.Guess, "minmax:", minmax_len)
		}
	}

	s.Guess = best.Guess
	s.Patterns = best.Patterns

	for _, p := range s.Patterns {
		// fmt.Println("Next Stage")
		BuildStrategy(p.NextStage, allwords)
	}
}

// XXX May be more efficient ways to do this, but let's just get started
// Does word match guess according to pattern?
func PatternMatch(guess string, word string, pattern []int) bool {

	for ip, cp := range pattern {
		cg := byte(guess[ip])
		// fmt.Println("ip", ip, "cp", cp, "cg", string(cg), "guess", guess, "word", word)

		// XXX clean this upto not use a switch statement.

		// If a 2, it needs to match it's a  matches, it's a 2
		if cp == 2 {
			if guess[ip] != word[ip] {
				// fmt.Println("false")
				return false
			}
			continue
		}

		// In the case of a 0 or 1, if it matches, it's a 2
		if (guess[ip] == word[ip]) {
			return false 
		}

		// We need to remove the perfect matches.  Then if there are still
		// some letters left over, it's a 1.  Otherwise, it's a 0.
		// XXX Is there a more efficient way to do this?

		// First eliminate dups by just putting in a wildcard char

		tmp := []byte(word)

		// Remove all the dups
		for idx, _ := range word {
			if cg == word[idx] && word[idx] == guess[idx] {
				tmp[idx] = '*'
			}
		}
		// fmt.Println("tmp", string(tmp))

		// Check if there are any remaining matching letters
		// If there are, it's a 1, it not, it's 0
		if cp == 0 {
			if str.IndexByte(string(tmp), cg) >= 0 {
				// fmt.Println("false")
				return false
			}
			continue
		}
		
		if cp == 1 {
			if str.IndexByte(string(tmp), cg) < 0 {
				// fmt.Println("false")
				return false
			}
			continue
		}
	}
	// fmt.Println("true")
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
