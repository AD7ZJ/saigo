package corpus

import (
	"sort"
	"strings"
	"unicode"
)

type tKeyValPair struct {
	Word  string
	Count int
}

func Analysis(textContent string) []tKeyValPair {
	histogramMap := make(map[string]int)
	var histogramSlice []tKeyValPair

	// Split the text into words
	words := splitIntoWords(textContent)

	// Count each word in the histogram
	for _, word := range words {
		histogramMap[word]++
	}

	// convert the resulting map into a slice
	for word, count := range histogramMap {
		histogramSlice = append(histogramSlice, tKeyValPair{Word: word, Count: count})
	}

	// Sort the slice by value (frequency of word occurence)
	sort.Slice(histogramSlice, func(i, j int) bool {
		return histogramSlice[i].Count > histogramSlice[j].Count // Sort in descending order
	})

	return histogramSlice
}

// splits a string into words, filtering out punctuation and converting to lowercase
func splitIntoWords(text string) []string {
	// Use FieldsFunc to split at any non-letter character
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c)
	})

	// Convert words to lowercase
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return words
}
