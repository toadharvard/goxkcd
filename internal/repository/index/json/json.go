package json

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/renameio"
	"github.com/toadharvard/goxkcd/internal/entity"
)

type JSONRepo struct {
	filePath string
}

func New(filePath string) (repo *JSONRepo) {
	return &JSONRepo{filePath: filePath}
}

func (r *JSONRepo) CreateOrUpdate(i entity.Index) error {
	file, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := json.Marshal(i)
	err = renameio.WriteFile(r.filePath, bytes, 0644)
	return err
}

func (r *JSONRepo) GetIndex() (entity.Index, error) {
	i := NewIndex()
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&i)
	return i, err
}

func (r *JSONRepo) Exists() bool {
	_, err := os.Stat(r.filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func (r *JSONRepo) BuildFromComics(comics []entity.Comix) (entity.Index, error) {
	index := NewIndex()
	for _, comic := range comics {
		for _, keyword := range comic.Keywords {
			index.Add(keyword, comic.ID)
		}
	}
	return index, nil
}
