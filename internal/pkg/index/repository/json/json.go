package json

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/renameio"
	"github.com/toadharvard/goxkcd/internal/pkg/index"
)

type IndexRepository struct {
	filePath string
}

func New(filePath string) (repo *IndexRepository) {
	repo = &IndexRepository{filePath: filePath}
	return
}

func (r *IndexRepository) CreateOrUpdate(i index.Index) error {
	file, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := json.Marshal(i)
	err = renameio.WriteFile(r.filePath, bytes, 0644)
	return err
}

func (r *IndexRepository) GetIndex() (index.Index, error) {
	i := index.New()
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&i)
	return i, err
}

func (r *IndexRepository) Exists() bool {
	_, err := os.Stat(r.filePath)
	return !errors.Is(err, os.ErrNotExist)
}
