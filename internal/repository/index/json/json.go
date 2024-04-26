package indexer

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/renameio"
	"github.com/toadharvard/goxkcd/internal/entity"
)

type JsonRepo struct {
	filePath string
}

func New(filePath string) (repo *JsonRepo) {
	return &JsonRepo{filePath: filePath}
}

func (r *JsonRepo) CreateOrUpdate(i *entity.Index) error {
	file, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := json.Marshal(i)
	err = renameio.WriteFile(r.filePath, bytes, 0644)
	return err
}

func (r *JsonRepo) GetIndex() (*entity.Index, error) {
	i := entity.NewIndex()
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&i)
	return i, err
}

func (r *JsonRepo) Exists() bool {
	_, err := os.Stat(r.filePath)
	return !errors.Is(err, os.ErrNotExist)
}
