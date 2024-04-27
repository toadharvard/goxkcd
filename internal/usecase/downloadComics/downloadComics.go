package downloadComics

import (
	"context"
	"sync"

	"github.com/toadharvard/goxkcd/internal/entity"
)

type ComixRepo interface {
	BulkInsert([]entity.Comix) error
	GetAll() ([]entity.Comix, error)
}

type XKCDRepo interface {
	GetByID(ctx context.Context, id int) (*entity.Comix, error)
	GetLastComixID(ctx context.Context) (int, error)
}

type UseCase struct {
	NumberOfWorkers int
	BatchSize       int
	comixRepo       ComixRepo
	xkcdRepo        XKCDRepo
}

func New(
	NumberOfWorkers int,
	BatchSize int,
	xkcdRepo XKCDRepo,
	comixRepo ComixRepo,
) *UseCase {
	return &UseCase{
		NumberOfWorkers: NumberOfWorkers,
		BatchSize:       BatchSize,
		xkcdRepo:        xkcdRepo,
		comixRepo:       comixRepo,
	}
}

func (u *UseCase) Run(ctx context.Context) (err error) {
	missingComixIDs := make(chan int)
	batches := make(chan []entity.Comix)
	go func() {
		err = u.writeMissingIDs(ctx, missingComixIDs)
	}()

	wg := sync.WaitGroup{}
	wg.Add(u.NumberOfWorkers)

	for i := 0; i < u.NumberOfWorkers; i++ {
		go func() {
			u.fetchComicsBatch(ctx, missingComixIDs, batches)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(batches)
	}()

	for batch := range batches {
		err = u.comixRepo.BulkInsert(batch)
		if err != nil {
			return
		}
	}
	return nil
}

func (u *UseCase) writeMissingIDs(ctx context.Context, missingComixIDs chan<- int) (err error) {
	defer close(missingComixIDs)
	limit, err := u.xkcdRepo.GetLastComixID(ctx)
	if err != nil {
		return err
	}

	allComics, err := u.comixRepo.GetAll()
	if err != nil {
		return err
	}

	exist := map[int]bool{}
	for _, comix := range allComics {
		exist[comix.ID] = true
	}

	for i := 1; i <= limit; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if !exist[i] {
				missingComixIDs <- i
			}
		}
	}
	return nil
}

func (u *UseCase) fetchComicsBatch(
	ctx context.Context,
	IDsToFetch <-chan int,
	batches chan<- []entity.Comix,
) {
	batch := []entity.Comix{}
	sendBatch := func() {
		if len(batch) > 0 {
			batches <- batch
			batch = []entity.Comix{}
		}
	}
	defer sendBatch()

	for {
		select {
		case id, ok := <-IDsToFetch:
			if !ok {
				return
			}

			comix, err := u.xkcdRepo.GetByID(ctx, id)
			if err == nil {
				batch = append(batch, *comix)
			}

			if len(batch) == u.BatchSize {
				sendBatch()
			}
		case <-ctx.Done():
			return
		}
	}
}
