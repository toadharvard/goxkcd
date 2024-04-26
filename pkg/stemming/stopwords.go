package stemming

import "github.com/toadharvard/goxkcd/pkg/iso6391"

func (s *stemmer) removeStopwords(tokens []Token, language iso6391.ISOCode6391) ([]Token, error) {
	tokensWithoutStopwords := []Token{}
	for _, token := range tokens {
		if !s.stopwords.IsStopword(token, language.Code) {
			tokensWithoutStopwords = append(tokensWithoutStopwords, token)
		}
	}
	return tokensWithoutStopwords, nil
}
