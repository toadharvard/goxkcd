package stemming

import (
	"strings"

	"github.com/kljensen/snowball"
	"github.com/toadharvard/goxkcd/internal/config"
)

func Stem(tokens []Token, language config.ISOCode639_1) ([]Token, error) {
	stemmedTokens := []Token{}

	snowballLang := getSnowballLanguageFromISOCode639_1(language)

	for _, Token := range tokens {
		stemmed, err := snowball.Stem(Token, snowballLang, false)
		if err != nil {
			return stemmedTokens, err
		}
		stemmedTokens = append(stemmedTokens, stemmed)
	}
	return stemmedTokens, nil
}

func StemString(str string, language config.ISOCode639_1) []Token {
	tokens := Tokenize(str)
	withoutDuplicates := RemoveDuplicates(tokens)
	withoutStopwords, _ := RemoveStopwords(withoutDuplicates, language)
	stemmedTokens, _ := Stem(withoutStopwords, language)
	return stemmedTokens
}

func getSnowballLanguageFromISOCode639_1(language config.ISOCode639_1) (newLang string) {
	switch strings.ToUpper(language) {
	case "EN":
		newLang = "english"
	case "ES":
		newLang = "spanish"
	case "FR":
		newLang = "french"
	case "RU":
		newLang = "russian"
	case "SV":
		newLang = "swedish"
	case "NO":
		newLang = "norwegian"
	case "HU":
		newLang = "hungarian"
	default:
		newLang = "english"
		return
	}
	return
}
