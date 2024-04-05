package repository

import (
	"encoding/json"
	"errors"
	"os"

	c "github.com/toadharvard/goxkcd/internal/pkg/comix"
)

type ComixRepository struct {
	filePath string
}

func New(filePath string) *ComixRepository {
	file, err := os.Create(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode([]c.Comix{}); err != nil {
		return nil
	}
	return &ComixRepository{filePath: filePath}
}

func (r *ComixRepository) GetAll() ([]c.Comix, error) {
	comics := []c.Comix{}
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&comics); err != nil {
		return nil, err
	}
	return comics, nil
}

func (r *ComixRepository) GetById(id int) (c.Comix, error) {
	comics, err := r.GetAll()
	if err != nil {
		return c.Comix{}, err
	}
	for _, comix := range comics {
		if comix.Id == id {
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
	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(comics)
}

func (r *ComixRepository) Exists() bool {
	_, err := os.Stat(r.filePath)
	return os.IsExist(err)
}
