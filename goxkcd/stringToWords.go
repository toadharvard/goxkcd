package main

import "strings"

// StringToWords splits a string into words and returns a slice of unique words.
//
// str: the input string to be split into words.
// []string: a slice of unique words extracted from the input string.
func StringToWords(str string) []string {
	wordsMap := make(map[string]bool)
	words := strings.Fields(str)
	var uniqueWords []string
	for _, word := range words {
		if !wordsMap[word] {
			uniqueWords = append(uniqueWords, word)
			wordsMap[word] = true
		}
	}
	return uniqueWords
}
