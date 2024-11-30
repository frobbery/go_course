package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(s string) []string {
	words := strings.Fields(s)
	countMap := make(map[string]int)
	for _, word := range words {
		countMap[word]++
	}
	uniqueWords := make([]string, 0, len(countMap))
	for word := range countMap {
		uniqueWords = append(uniqueWords, word)
	}
	sort.Slice(uniqueWords, func(i, j int) bool {
		if countMap[uniqueWords[i]] == countMap[uniqueWords[j]] {
			return uniqueWords[i] < uniqueWords[j]
		}
		return countMap[uniqueWords[i]] > countMap[uniqueWords[j]]
	})
	topLen := 10
	if cap(uniqueWords) < topLen {
		topLen = cap(uniqueWords)
	}
	return uniqueWords[:topLen]
}
