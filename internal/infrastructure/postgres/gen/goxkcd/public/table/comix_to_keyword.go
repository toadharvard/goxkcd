//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var ComixToKeyword = newComixToKeywordTable("public", "comix_to_keyword", "")

type comixToKeywordTable struct {
	postgres.Table

	// Columns
	ComixID   postgres.ColumnInteger
	KeywordID postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ComixToKeywordTable struct {
	comixToKeywordTable

	EXCLUDED comixToKeywordTable
}

// AS creates new ComixToKeywordTable with assigned alias
func (a ComixToKeywordTable) AS(alias string) *ComixToKeywordTable {
	return newComixToKeywordTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ComixToKeywordTable with assigned schema name
func (a ComixToKeywordTable) FromSchema(schemaName string) *ComixToKeywordTable {
	return newComixToKeywordTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ComixToKeywordTable with assigned table prefix
func (a ComixToKeywordTable) WithPrefix(prefix string) *ComixToKeywordTable {
	return newComixToKeywordTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ComixToKeywordTable with assigned table suffix
func (a ComixToKeywordTable) WithSuffix(suffix string) *ComixToKeywordTable {
	return newComixToKeywordTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newComixToKeywordTable(schemaName, tableName, alias string) *ComixToKeywordTable {
	return &ComixToKeywordTable{
		comixToKeywordTable: newComixToKeywordTableImpl(schemaName, tableName, alias),
		EXCLUDED:            newComixToKeywordTableImpl("", "excluded", ""),
	}
}

func newComixToKeywordTableImpl(schemaName, tableName, alias string) comixToKeywordTable {
	var (
		ComixIDColumn   = postgres.IntegerColumn("comix_id")
		KeywordIDColumn = postgres.IntegerColumn("keyword_id")
		allColumns      = postgres.ColumnList{ComixIDColumn, KeywordIDColumn}
		mutableColumns  = postgres.ColumnList{}
	)

	return comixToKeywordTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ComixID:   ComixIDColumn,
		KeywordID: KeywordIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
