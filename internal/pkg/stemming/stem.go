package stemming

import (
	"strings"

	"github.com/kljensen/snowball"
	"github.com/toadharvard/goxkcd/internal/config"
	sw "github.com/toadharvard/stopwords-iso"
)

type Stemmer struct {
	stopwords sw.StopwordsMapping
}

func New() *Stemmer {
	stopwords, _ := sw.NewStopwordsMapping()
	stopwords["en"] = append(stopwords["en"], "alt")
	stopwords["en"] = append(stopwords["en"], "text")
	return &Stemmer{stopwords: stopwords}
}

func (s *Stemmer) Stem(tokens []Token, language config.ISOCode639_1) ([]Token, error) {
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

func (s *Stemmer) StemString(str string, language config.ISOCode639_1) []Token {
	tokens := Tokenize(str)
	withoutDuplicates := RemoveDuplicates(tokens)
	withoutStopwords, _ := s.RemoveStopwords(withoutDuplicates, language)
	stemmedTokens, _ := s.Stem(withoutStopwords, language)
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
