package stemming

import (
	"sync"

	"github.com/kljensen/snowball"
	"github.com/toadharvard/goxkcd/pkg/iso6391"
	sw "github.com/toadharvard/stopwords-iso"
)

type stemmer struct {
	stopwords sw.StopwordsMapping
}

var onceStemmer sync.Once
var stemmerInstance *stemmer

func New() *stemmer {
	onceStemmer.Do(func() {
		stopwords, _ := sw.NewStopwordsMapping()
		stopwords["en"] = append(stopwords["en"], "alt")
		stopwords["en"] = append(stopwords["en"], "text")
		stopwords["en"] = append(stopwords["en"], "title")
		stemmerInstance = &stemmer{stopwords: stopwords}
	})
	return stemmerInstance
}

func (s *stemmer) Stem(tokens []Token, language iso6391.ISOCode6391) ([]Token, error) {
	stemmedTokens := []Token{}

	for _, Token := range tokens {
		stemmed, err := snowball.Stem(Token, language.Name, false)
		if err != nil {
			return stemmedTokens, err
		}
		stemmedTokens = append(stemmedTokens, stemmed)
	}
	return stemmedTokens, nil
}

func (s *stemmer) StemString(str string, language iso6391.ISOCode6391) []Token {
	tokens := Tokenize(str)
	withoutDuplicates := RemoveDuplicates(tokens)
	withoutStopwords, _ := s.removeStopwords(withoutDuplicates, language)
	stemmedTokens, _ := s.Stem(withoutStopwords, language)
	return stemmedTokens
}
