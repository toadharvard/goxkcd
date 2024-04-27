package suggestComix

import (
	"log/slog"

	"github.com/toadharvard/goxkcd/internal/entity"
	"github.com/toadharvard/goxkcd/pkg/iso6391"
	"github.com/toadharvard/goxkcd/pkg/stemming"
)

type IndexRepo interface {
	GetIndex() (*entity.Index, error)
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
	suggestionIDs := index.GetRelevantIDs(keywords)
	limit := min(suggestionsLimit, len(suggestionIDs))
	slog.Info("number of suggestions", "limit", limit)
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
