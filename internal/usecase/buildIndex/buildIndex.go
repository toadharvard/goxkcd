package buildIndex

import (
	"context"

	"github.com/toadharvard/goxkcd/internal/entity"
)

type IndexRepo interface {
	CreateOrUpdate(*entity.Index) error
}

type ComixRepo interface {
	GetAll() ([]entity.Comix, error)
}

type BuildIndexUC struct {
	indexRepo IndexRepo
	comixRepo ComixRepo
}

func New(indexRepo IndexRepo, comixRepo ComixRepo) *BuildIndexUC {
	return &BuildIndexUC{
		indexRepo: indexRepo,
		comixRepo: comixRepo,
	}
}

func (u *BuildIndexUC) Run(ctx context.Context) (err error) {
	comics, err := u.comixRepo.GetAll()
	if err != nil {
		return
	}

	index := entity.NewIndexFromComics(comics)

	err = u.indexRepo.CreateOrUpdate(index)
	return
}
