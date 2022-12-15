package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type word struct {
	count int
	value string
}

type parser struct {
	pattern    string
	rowData    string
	parsedData map[string]int
	words      []word
}

func (p *parser) parse() {
	reg, _ := regexp.Compile(p.pattern)

	for _, word := range reg.FindAllString(p.rowData, -1) {
		p.parsedData[strings.ToLower(word)]++
	}

	p.words = make([]word, 0, len(p.parsedData))

	for value, count := range p.parsedData {
		p.words = append(p.words, word{count, value})
	}
}

func (p *parser) sort() {
	if len(p.words) == 0 {
		return
	}

	sort.Slice(
		p.words,
		func(i, j int) bool {
			if p.words[i].count == p.words[j].count {
				return p.words[i].value < p.words[j].value
			}

			return p.words[i].count > p.words[j].count
		},
	)
}

func (p *parser) slice(n int) []string {
	if len(p.words) == 0 {
		return nil
	}

	if n > len(p.words) {
		n = len(p.words)
	}

	result := make([]string, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, p.words[i].value)
	}

	return result
}

func Top10(str string) []string {
	pattern := `\S+`

	p := parser{
		pattern:    pattern,
		rowData:    str,
		parsedData: make(map[string]int),
	}

	p.parse()
	p.sort()

	return p.slice(10)
}
