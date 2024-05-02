package xkcd

import (
	"context"

	"github.com/toadharvard/goxkcd/internal/entity"
)

type XKCDClient interface {
	GetComixByID(ctx context.Context, id int) (entity.Comixer, error)
	GetLastComix(ctx context.Context) (entity.Comixer, error)
}

type XKCDComixRepo struct {
	client XKCDClient
}

func New(client XKCDClient) *XKCDComixRepo {
	return &XKCDComixRepo{
		client: client,
	}
}

func (r *XKCDComixRepo) GetByID(ctx context.Context, id int) (*entity.Comix, error) {
	comix, err := r.client.GetComixByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return comix.ToComixEntity(), nil
}

func (r *XKCDComixRepo) GetLastComix(ctx context.Context) (*entity.Comix, error) {
	comix, err := r.client.GetLastComix(ctx)
	return comix.ToComixEntity(), err
}
