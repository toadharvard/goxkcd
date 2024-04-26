package xkcd

import (
	"context"

	"github.com/toadharvard/goxkcd/internal/entity"
)

type XKCDComixRepo struct {
	client *XKCDClient
}

func New(client *XKCDClient) *XKCDComixRepo {
	return &XKCDComixRepo{
		client: client,
	}
}

func (r *XKCDComixRepo) GetByID(ctx context.Context, id int) (*entity.Comix, error) {
	comix, err := r.client.GetComixByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return comix.ToComix(), nil
}

func (r *XKCDComixRepo) GetLastComixID(ctx context.Context) (id int, err error) {
	return r.client.GetLastComixID(ctx)
}
