package buildindex

import (
	"context"

	"github.com/toadharvard/goxkcd/internal/entity"
)

type IndexRepo interface {
	BuildFromComics([]entity.Comix) (entity.Index, error)
	CreateOrUpdate(entity.Index) error
}

type ComixRepo interface {
	GetAll() ([]entity.Comix, error)
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

func (u *UseCase) Run(ctx context.Context) (err error) {
	comics, err := u.comixRepo.GetAll()
	if err != nil {
		return
	}

	index, err := u.indexRepo.BuildFromComics(comics)
	if err != nil {
		return
	}
	err = u.indexRepo.CreateOrUpdate(index)
	return
}
