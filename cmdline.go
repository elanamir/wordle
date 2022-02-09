package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
  "encoding/json"
	"os"
	str "strings"
)

func LaunchTool(sfile string, wordlen int) {

	cur := ReadStrategy(sfile, wordlen)


	alltwos := make([]int, wordlen)
	for i, _ := range alltwos {
		alltwos[i] = 2
	}

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("Guess: ", cur.Guess)
    if len(cur.Next) == 0 {
      fmt.Println("All done!")
      os.Exit(0)
    }
		fmt.Print("Response: ")
		reponseString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		reponseString = str.TrimSuffix(reponseString, "\n")

		pattern := String2Pattern(reponseString)

		if EqualPattern(pattern, alltwos) {
			fmt.Println("Hooray!")
			break
		}

		nopattern := true
		for _, rec := range cur.Next {
			if EqualPattern(rec.Pattern, pattern) {
				cur = rec
				nopattern = false
				break
			}
		}
		if nopattern {
			panic("No pattern!")
		}
	}
}

func ReadStrategy(sfile string, wordlen int) JsonRec {
	jsonFile, err := os.Open(sfile)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result JsonRec
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

func EqualPattern(p1 Pattern, p2 Pattern) bool {
  for i, p := range p1 {
    if p != p2[i] {
      return false
    }
  }
  return true
}

func String2Pattern(s string) Pattern {
	ret := []int{}
	for _, c := range s {
		switch c {
		case '0':
			ret = append(ret, 0)
		case '1':
			ret = append(ret, 1)
		case '2':
			ret = append(ret, 2)
		case ' ':
		default:
			return []int{}
		}
	}
	return ret
}
