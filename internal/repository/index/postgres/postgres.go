package postgres

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres" //lint:ignore ST1001 pretty sql-like syntax
	_ "github.com/lib/pq"
	"github.com/toadharvard/goxkcd/internal/entity"
	"github.com/toadharvard/goxkcd/internal/infrastructure/postgres/gen/goxkcd/public/table"
)

type PostgresIndex struct {
	dsn string
	db  *sql.DB
}

func NewIndex(dsn string) (*PostgresIndex, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresIndex{
		dsn: dsn,
		db:  db,
	}, nil
}
func (i *PostgresIndex) Search(token entity.Token) []int {
	searchQuery := SELECT(
		table.Keyword.Word,
		table.ComixToKeyword.ComixID,
	).FROM(
		table.ComixToKeyword.INNER_JOIN(
			table.Keyword,
			table.Keyword.ID.EQ(table.ComixToKeyword.KeywordID),
		),
	).WHERE(table.Keyword.Word.EQ(String(token)))

	var comixIDs []struct {
		ID int `alias:"comix_to_keyword.comix_id"`
	}

	err := searchQuery.Query(i.db, &comixIDs)
	if err != nil {
		return nil
	}

	ids := make([]int, len(comixIDs))
	for i, id := range comixIDs {
		ids[i] = id.ID
	}
	return ids
}

type PostgresRepo struct {
	db  *sql.DB
	dsn string
}

func New(dsn string) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{
		db:  db,
		dsn: dsn,
	}, nil
}

func (r *PostgresRepo) GetIndex() (entity.Index, error) {
	return NewIndex(r.dsn)
}

func (r *PostgresRepo) CreateOrUpdate(i entity.Index) error {
	// This method is a no-op, as the database will handle index updates automatically
	// upon insertion of new Comix rows.
	return nil
}

func (r *PostgresRepo) BuildFromComics([]entity.Comix) (entity.Index, error) {
	// This method is a no-op, as the database will handle index updates automatically
	// upon insertion of new Comix rows.
	return NewIndex(r.dsn)
}
