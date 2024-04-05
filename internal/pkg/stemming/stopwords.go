package stemming

import (
	"github.com/toadharvard/goxkcd/internal/config"
	sw "github.com/toadharvard/stopwords-iso"
)

func RemoveStopwords(tokens []Token, language config.ISOCode639_1) ([]Token, error) {
	m, err := sw.NewStopwordsMapping()
	if err != nil {
		return tokens, err
	}

	tokensWithoutStopwords := []Token{}

	for _, Token := range tokens {
		if !m.IsStopword(Token, language) {
			tokensWithoutStopwords = append(tokensWithoutStopwords, Token)
		}
	}

	return tokensWithoutStopwords, nil
}
