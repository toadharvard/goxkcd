package suggestcomix

import (
	"log/slog"
	"sort"

	"github.com/toadharvard/goxkcd/internal/entity"
	"github.com/toadharvard/goxkcd/pkg/iso6391"
	"github.com/toadharvard/goxkcd/pkg/stemming"
)

type IndexRepo interface {
	GetIndex() (entity.Index, error)
}

type ComixRepo interface {
	GetByID(int) (entity.Comix, error)
}

type UseCase struct {
	indexRepo IndexRepo
	comixRepo ComixRepo
}

func New(indexRepo IndexRepo, comixRepo ComixRepo) *UseCase {
	return &UseCase{
		indexRepo: indexRepo,
		comixRepo: comixRepo,
	}
}

func (u *UseCase) Run(
	language iso6391.ISOCode6391,
	searchQuery string,
	suggestionsLimit int,
) ([]entity.Comix, error) {
	keywords := stemming.New().StemString(searchQuery, language)
	index, err := u.indexRepo.GetIndex()
	if err != nil {
		return nil, err
	}

	slog.Info("Suggesting", "keywords", keywords, "limit", suggestionsLimit, "query", searchQuery)

	suggestionIDs := getRelevantIDs(index, keywords)
	limit := min(suggestionsLimit, len(suggestionIDs))
	suggestions := make([]entity.Comix, 0, limit)
	for _, id := range suggestionIDs[:limit] {
		comix, err := u.comixRepo.GetByID(id)
		if err != nil {
			return nil, err
		}
		suggestions = append(suggestions, comix)
	}
	return suggestions, nil
}

func getRelevantIDs(i entity.Index, tokens []entity.Token) []int {
	counter := make(map[int]int)
	for _, token := range tokens {
		for _, id := range i.Search(token) {
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
