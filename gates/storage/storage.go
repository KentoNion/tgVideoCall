package storage

import (
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/bool64/sqluct"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db  *sqlx.DB
	sq  sq.StatementBuilderType
	sm  sqluct.Mapper
	log slog.Logger
}

func NewPostgresDB(db *sqlx.DB, log slog.Logger) *DB {
	dab := DB{
		db:  db,
		sm:  sqluct.Mapper{Dialect: sqluct.DialectPostgres},
		sq:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		log: log,
	}
	_, err := dab.db.Exec("SET timezone TO 'UTC'")
	if err != nil {
		panic(err)
	}
	return &dab
}
