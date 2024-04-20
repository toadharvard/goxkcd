package index

import (
	"sort"

	"github.com/toadharvard/goxkcd/internal/pkg/comix"
)

type ID = int
type Token = string
type Index map[Token][]ID

func New() Index {
	return Index{}
}

func (i Index) Add(token Token, id ID) {
	i[token] = append(i[token], id)
}

func FromComics(comics []comix.Comix) Index {
	index := New()
	for _, comic := range comics {
		for _, keyword := range comic.Keywords {
			index.Add(keyword, comic.ID)
		}
	}
	return index
}

func (i Index) GetRelevantIDs(tokens []Token) []ID {
	counter := make(map[ID]int)
	for _, token := range tokens {
		for _, id := range i[token] {
			counter[id]++
		}
	}

	type pair struct {
		id          ID
		occurrences int
	}

	var relevances []pair
	for id, count := range counter {
		relevances = append(relevances, pair{id, count})
	}

	sort.Slice(relevances, func(i, j int) bool {
		return relevances[i].occurrences > relevances[j].occurrences
	})

	ids := make([]ID, 0, len(relevances))
	for _, pair := range relevances {
		ids = append(ids, pair.id)
	}

	return ids
}
