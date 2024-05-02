package json

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"

	"github.com/google/renameio"
	"github.com/toadharvard/goxkcd/internal/entity"
)

var permissions fs.FileMode = 0644

type JSONRepo struct {
	filePath string
}

func New(filePath string) (repo *JSONRepo) {
	repo = &JSONRepo{filePath: filePath}
	return
}

func (r *JSONRepo) Create() error {
	file, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := json.Marshal([]entity.Comix{})
	err = renameio.WriteFile(r.filePath, bytes, permissions)
	return err
}

func (r *JSONRepo) GetAll() ([]entity.Comix, error) {
	comics := []entity.Comix{}
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&comics)
	return comics, err
}

func (r *JSONRepo) GetByID(id int) (entity.Comix, error) {
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

func (r *JSONRepo) BulkInsert(toInsert []entity.Comix) error {
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

func (r *JSONRepo) Exists() bool {
	_, err := os.Stat(r.filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func (r *JSONRepo) Size() (int, error) {
	comics, err := r.GetAll()
	if err != nil {
		return 0, err
	}
	return len(comics), nil
}
