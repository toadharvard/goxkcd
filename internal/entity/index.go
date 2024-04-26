package entity

import (
	"sort"
)

type Token = string
type Index struct {
	Values map[Token][]int `json:"values"`
}

func NewIndex() *Index {
	values := make(map[Token][]int)
	return &Index{
		Values: values,
	}
}

func NewIndexFromComics(comics []Comix) *Index {
	index := NewIndex()
	for _, comic := range comics {
		for _, keyword := range comic.Keywords {
			index.Add(keyword, comic.ID)
		}
	}
	return index
}

func (i *Index) Add(token Token, id int) {
	i.Values[token] = append(i.Values[token], id)
}

func (i Index) GetRelevantIDs(tokens []Token) []int {
	counter := make(map[int]int)
	for _, token := range tokens {
		for _, id := range i.Values[token] {
			counter[id]++
		}
	}

	type pair struct {
		id          int
		occurrences int
	}

	var relevances []pair
	for id, count := range counter {
		relevances = append(relevances, pair{id, count})
	}

	sort.Slice(relevances, func(i, j int) bool {
		return relevances[i].occurrences > relevances[j].occurrences
	})

	ids := make([]int, 0, len(relevances))
	for _, pair := range relevances {
		ids = append(ids, pair.id)
	}

	return ids
}
