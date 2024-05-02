package json

import "github.com/toadharvard/goxkcd/internal/entity"

type IndexJSON struct {
	Values map[entity.Token][]int `json:"values"`
}

func NewIndex() *IndexJSON {
	values := make(map[entity.Token][]int)
	return &IndexJSON{
		Values: values,
	}
}

func (i *IndexJSON) Add(token entity.Token, id int) {
	i.Values[token] = append(i.Values[token], id)
}

func (i *IndexJSON) Search(token entity.Token) []int {
	return i.Values[token]
}
