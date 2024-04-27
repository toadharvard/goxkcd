package comix

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"

	"github.com/google/renameio"
	"github.com/toadharvard/goxkcd/internal/entity"
)

var permissions fs.FileMode = 0644

type JsonRepo struct {
	filePath string
}

func New(filePath string) (repo *JsonRepo) {
	repo = &JsonRepo{filePath: filePath}
	return
}

func (r *JsonRepo) Create() error {
	file, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := json.Marshal([]entity.Comix{})
	err = renameio.WriteFile(r.filePath, bytes, permissions)
	return err
}

func (r *JsonRepo) GetAll() ([]entity.Comix, error) {
	comics := []entity.Comix{}
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&comics)
	return comics, err
}

func (r *JsonRepo) GetByID(id int) (entity.Comix, error) {
	comics, err := r.GetAll()
	if err != nil {
		return entity.Comix{}, err
	}
	for _, comix := range comics {
		if comix.ID == id {
			return comix, nil
		}
	}
	return entity.Comix{}, errors.New("comix not found")
}

func (r *JsonRepo) BulkInsert(toInsert []entity.Comix) error {
	comics, err := r.GetAll()
	if err != nil {
		return err
	}
	comics = append(comics, toInsert...)
	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, permissions)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := json.Marshal(comics)
	return renameio.WriteFile(r.filePath, bytes, permissions)
}

func (r *JsonRepo) Exists() bool {
	_, err := os.Stat(r.filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func (r *JsonRepo) Size() (int, error) {
	comics, err := r.GetAll()
	if err != nil {
		return 0, err
	}
	return len(comics), nil
}
