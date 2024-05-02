package postgres

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres" //lint:ignore ST1001 pretty sql-like syntax
	_ "github.com/lib/pq"
	"github.com/toadharvard/goxkcd/internal/entity"
	"github.com/toadharvard/goxkcd/internal/infrastructure/postgres/gen/goxkcd/public/model"
	"github.com/toadharvard/goxkcd/internal/infrastructure/postgres/gen/goxkcd/public/table"
)

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

func (r *PostgresRepo) GetAll() ([]entity.Comix, error) {
	var comixModels []struct {
		model.Comix
		Keywords []string `alias:"keyword.word"`
	}

	query := SELECT(
		table.Comix.AllColumns,
		table.Keyword.Word,
	).FROM(
		table.Comix.INNER_JOIN(table.ComixToKeyword, table.Comix.ID.EQ(table.ComixToKeyword.ComixID)).
			INNER_JOIN(table.Keyword, table.Keyword.ID.EQ(table.ComixToKeyword.KeywordID)),
	)

	err := query.Query(r.db, &comixModels)
	if err != nil {
		return nil, err
	}

	comixes := make([]entity.Comix, len(comixModels))
	for i, comixModel := range comixModels {
		comixes[i] = *entity.NewComix(int(comixModel.ID), comixModel.URL, comixModel.Keywords)
	}

	return comixes, nil
}

func (r *PostgresRepo) GetByID(id int) (entity.Comix, error) {
	ID := Int(int64(id))
	var comixModel struct {
		model.Comix
		Keywords []string `alias:"keyword.word"`
	}

	query := SELECT(
		table.Comix.AllColumns,
		table.Keyword.Word,
	).FROM(
		table.Comix.INNER_JOIN(table.ComixToKeyword, table.Comix.ID.EQ(table.ComixToKeyword.ComixID)).
			INNER_JOIN(table.Keyword, table.Keyword.ID.EQ(table.ComixToKeyword.KeywordID)),
	).WHERE(
		table.Comix.ID.EQ(ID),
	)

	err := query.Query(r.db, &comixModel)
	if err != nil {
		return entity.Comix{}, err
	}

	comix := entity.NewComix(int(comixModel.ID), comixModel.URL, comixModel.Keywords)
	return *comix, nil
}

func (r *PostgresRepo) BulkInsert(toInsert []entity.Comix) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, comix := range toInsert {
		insertComixQuery := table.Comix.INSERT(
			table.Comix.AllColumns,
		).MODEL(model.Comix{
			ID:  int32(comix.ID),
			URL: comix.URL,
		}).ON_CONFLICT(table.Comix.ID).DO_UPDATE(
			SET(
				table.Comix.URL.SET(String(comix.URL)),
			),
		)

		_, err = insertComixQuery.Exec(tx)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		for _, keyword := range comix.Keywords {
			insertKeywordQuery := table.Keyword.INSERT(
				table.Keyword.Word,
			).MODEL(model.Keyword{
				Word: keyword,
			}).ON_CONFLICT(table.Keyword.Word).DO_NOTHING()

			_, err := insertKeywordQuery.Exec(tx)

			if err != nil {
				_ = tx.Rollback()
				return err
			}

			selectKeywordQuery := SELECT(table.Keyword.ID).FROM(table.Keyword).WHERE(table.Keyword.Word.EQ(String(keyword)))
			var keywordModel struct {
				ID int `alias:"keyword.id"`
			}
			err = selectKeywordQuery.Query(tx, &keywordModel)
			if err != nil {
				_ = tx.Rollback()
				return err
			}

			insertComixToKeywordQuery := table.ComixToKeyword.INSERT(
				table.ComixToKeyword.ComixID, table.ComixToKeyword.KeywordID,
			).MODEL(model.ComixToKeyword{
				ComixID:   int32(comix.ID),
				KeywordID: int32(keywordModel.ID),
			}).ON_CONFLICT().DO_NOTHING()

			_, err = insertComixToKeywordQuery.Exec(tx)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}
	err = tx.Commit()
	return err
}

func (r *PostgresRepo) Size() (int, error) {
	var model struct {
		Count int
	}
	query := SELECT(COUNT(table.Comix.ID)).FROM(table.Comix)
	err := query.Query(r.db, &model)
	return model.Count, err
}
