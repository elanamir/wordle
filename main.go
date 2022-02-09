package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

type Pattern []int

var AllPatterns []Pattern

type JsonRec struct {
	Pattern Pattern   `json:"pattern,omitempty"`
	Guess   string    `json:"guess"`
	Next    []JsonRec `json:"next,omitempty"`
}

func main() {

	fnamePtr := flag.String("wordfile", "./words/CollinsWords-5.txt", "wordlist file")
	nwordsPtr := flag.Int("nwords", -1, "number of words to process (all if -1)")
	wordlenPtr := flag.Int("wordlen", 5, "Length of words in word list")
	otypePtr := flag.String("otype", "json", "Output type [json or native]")
	withdictPtr := flag.Bool("withdict", false, "Include dictionaries in native output")
	sfilePtr := flag.String("strategyfile", "s.json", "Strategy json file input")
	cmdPtr := flag.Bool("cmdline", false, "Launch command line")

	flag.Parse()

	if *cmdPtr {
		LaunchTool(*sfilePtr, *wordlenPtr)
		os.Exit(0)
	}

	AllPatterns = make([]Pattern, 0)
	BuildPattern(0, Pattern{}, *wordlenPtr)

	allwords := BuildWordList(*fnamePtr, *nwordsPtr)
	root := InitStrategy(allwords)

	switch *otypePtr {
	case "native":
		maxdepth := -1
		fmt.Println("Num Words: ", len(allwords))
		PrintStrategy(&root, 0, nil, &maxdepth, *withdictPtr)
		fmt.Println("Max Depth: ", maxdepth)

	case "json":
		var rec JsonRec
		JsonPrintStrategy(&root, &rec, nil)
		res, err := json.Marshal(rec)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(res))

	default:
		panic("Unknown output type")
	}
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
	Score      float64
	Patterns   []StrategyElement
}

func InitStrategy(allwords []string) StrategyStage {
	root := StrategyStage{Dictionary: allwords}
	BuildStrategy(&root, allwords)

	return root
}

func PrintStrategy(s *StrategyStage, depth int, pattern Pattern, maxdepth *int, withdict bool) {
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
			if withdict {
				fmt.Print(" Dict: ", s.Dictionary)
			}
		}

		fmt.Println("")

		if len(s.Patterns) != 0 {
			for _, p := range s.Patterns {
				if p.NextStage != nil && len(p.NextStage.Dictionary) != 0 {
					PrintStrategy(p.NextStage, depth+1, p.Pattern, maxdepth, withdict)
				}
			}
		}
	}
}

func JsonPrintStrategy(s *StrategyStage, rec *JsonRec, pattern Pattern) {
	if s == nil {
		return
	}

	rec.Pattern = pattern

	if len(s.Dictionary) != 0 {
		if len(s.Dictionary) == 1 {
			rec.Guess = s.Dictionary[0]
		} else {
			rec.Guess = s.Guess
		}

		if len(s.Patterns) != 0 {
			for _, p := range s.Patterns {
				if p.NextStage != nil && len(p.NextStage.Dictionary) != 0 {
					var t JsonRec
					JsonPrintStrategy(p.NextStage, &t, p.Pattern)
					rec.Next = append(rec.Next, t)
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

func OneGuess(guess string, dictionary []string) StrategyStage {
	cur := StrategyStage{Guess: guess}

	// Determine how many words from the dictionary match for each pattern
	for _, pattern := range AllPatterns {
		se := StrategyElement{Pattern: pattern}
		se.NextStage = &StrategyStage{}

		for _, word := range dictionary {
			if word != guess && PatternMatch(guess, word, pattern) {

				// Add the word to the pattern's dictionary
				se.NextStage.Dictionary = append(se.NextStage.Dictionary, word)
			}
		}

		if len(se.NextStage.Dictionary) != 0 {
			cur.Patterns = append(cur.Patterns, se)
		}
	}
	cur.CalcScore()
	
	return cur
}

func BuildStrategy(s *StrategyStage, allwords []string) {

	// fmt.Println("Entered")
	if len(s.Dictionary) < 2 {
		// fmt.Println("returning")
		return
	}

	// fmt.Println("strategy stage", s)

	best := StrategyStage{Score: InitScore()}

	for _, guess := range allwords {
		cur := OneGuess(guess, s.Dictionary)
		if cur.IsBest(best) {
				best = cur
		}
	}
	
	s.Guess = best.Guess
	s.Patterns = best.Patterns

	for _, p := range s.Patterns {
		// fmt.Println("Next Stage")
		BuildStrategy(p.NextStage, allwords)
	}
}

func InitScore() float64 {
	return 100000.0 //XXX
}

func (cur StrategyStage) IsBest(best StrategyStage) bool {
	return best.Score > cur.Score
}

func (cur *StrategyStage) CalcScore() {
	var n float64 = 0
	var sum float64 = 0
	for _, p := range cur.Patterns {
		l := len(p.NextStage.Dictionary)
		if l != 0 {
			sum += float64(l)
			n += 1
		}
	}

	cur.Score = sum / n
}

func PatternMatch(guess string, word string, pattern Pattern) bool {
	// Calculate the pattern between guess and word and then compare to pattern
	// Not efficient, but not sure there's a much better way to do this and this is certainly simplest.

	p := make(Pattern, len(word))

	bword := []byte(word)
	bguess := []byte(guess)

	// Determine the 2's
	for ig, cg := range bguess {
		p[ig] = 0 // initialize explicitly for good measure
		if cg == bword[ig] {
			p[ig] = 2
			bword[ig] = '*'
			bguess[ig] = '*'
		}
	}

	for ig, cg := range bguess {
		if cg == '*' {
			continue
		}
		for iw, cw := range bword {
			if cw == '*' {
				continue
			}
			if cg == cw {
				p[ig] = 1
				bword[iw] = '*'
				break
			}
		}
	}

	// Compare the patterns
	for ip, cp := range pattern {
		if p[ip] != cp {
			return false
		}
	}
	return true
}
