package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func Top10(text string) []string {
	wordMap := make(map[string]uint)
	for _, w := range strings.Fields(text) {
		w = strings.TrimFunc(w, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		if len(w) != 0 {
			wordMap[strings.ToLower(w)]++
		}
	}

	res := make([]string, 0, len(wordMap))
	for key := range wordMap {
		res = append(res, key)
	}
	sort.Slice(res, func(i, j int) bool {
		if wordMap[res[i]] > wordMap[res[j]] ||
			(wordMap[res[i]] == wordMap[res[j]] && strings.Compare(res[i], res[j]) < 0) {
			return true
		}
		return false
	})
	maxLen := 10
	if len(res) < 10 {
		maxLen = len(res)
	}
	return res[0:maxLen]
}
