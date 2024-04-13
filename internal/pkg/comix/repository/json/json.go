package repository

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/renameio"
	c "github.com/toadharvard/goxkcd/internal/pkg/comix"
)

type ComixRepository struct {
	filePath string
}

func New(filePath string) (repo *ComixRepository, err error) {
	repo = &ComixRepository{filePath: filePath}
	if !repo.Exists() {
		err = repo.Create()
	}
	return
}

func (r *ComixRepository) Create() error {
	file, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := json.Marshal([]c.Comix{})
	err = renameio.WriteFile(r.filePath, bytes, 0644)
	return err
}

func (r *ComixRepository) GetAll() ([]c.Comix, error) {
	comics := []c.Comix{}
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&comics)
	return comics, err
}

func (r *ComixRepository) GetByID(id int) (c.Comix, error) {
	comics, err := r.GetAll()
	if err != nil {
		return c.Comix{}, err
	}
	for _, comix := range comics {
		if comix.ID == id {
			return comix, nil
		}
	}
	return c.Comix{}, errors.New("comix not found")
}

func (r *ComixRepository) BulkInsert(comixList []c.Comix) error {
	comics, err := r.GetAll()
	if err != nil {
		return err
	}
	comics = append(comics, comixList...)
	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := json.Marshal(comics)
	return renameio.WriteFile(r.filePath, bytes, 0644)
}

func (r *ComixRepository) Exists() bool {
	_, err := os.Stat(r.filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func (r *ComixRepository) Size() int {
	comics, err := r.GetAll()
	if err != nil {
		return 0
	}
	return len(comics)
}
