package stemming

import (
	"regexp"
	"strings"
)

type Token = string

var alphanumericOnly = regexp.MustCompile(`[^\p{L}\p{N} ]+`)
var wordSegmenter = regexp.MustCompile(`[\pL\p{Mc}\p{Mn}-_']+`)

func Tokenize(str string) []Token {
	words := strings.Fields(str)
	resplitted := strings.Join(words, " ")
	onlyAlphanumeric := alphanumericOnly.ReplaceAllString(resplitted, "")
	words = wordSegmenter.FindAllString(onlyAlphanumeric, -1)
	return words
}

func RemoveDuplicates(tokens []Token) []Token {
	keys := make(map[Token]bool)
	list := []Token{}
	for _, entry := range tokens {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
