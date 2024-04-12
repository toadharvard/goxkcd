package stemming

import (
	"github.com/toadharvard/goxkcd/internal/config"
)

func (s *stemmer) removeStopwords(tokens []Token, language config.ISOCode639_1) ([]Token, error) {
	tokensWithoutStopwords := []Token{}
	for _, token := range tokens {
		if !s.stopwords.IsStopword(token, language) {
			tokensWithoutStopwords = append(tokensWithoutStopwords, token)
		}
	}
	return tokensWithoutStopwords, nil
}
