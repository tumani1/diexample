package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"gitlab.com/igor.tumanov1/theboatscom/internal/domain"
)

type scanner interface {
	Scan(dest ...interface{}) error
}

type searchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) domain.ISearchRepository {
	return &searchRepository{
		db: db,
	}
}

func (r *searchRepository) Find(ctx context.Context, filter *domain.SearchFilter) (_ []*domain.SearchEntry, err error) {
	const query = `
WITH search_query AS (
    SELECT ym.id, ym.name model_name, yb.name builder_name, similarity(ym.name, $1) score
    FROM yacht_models ym
		JOIN yacht_builders yb ON ym.builder_id = yb.id
    WHERE ym.name ILIKE $1 || '%'

    UNION

    SELECT ym.id, ym.name model_name, yb.name builder_name, similarity(yb.name, $1) score
    FROM yacht_models ym
		JOIN yacht_builders yb ON ym.builder_id = yb.id
    WHERE yb.name ILIKE $1 || '%'
)

SELECT t2.id boat_id, t3.name fleet_name, t1.model_name, t1.builder_name
FROM search_query t1
	JOIN boats t2 ON t2.model_id = t1.id
	JOIN fleets t3 ON t2.fleet_id = t3.id
ORDER BY t1.score DESC
OFFSET $2 LIMIT $3
`

	var rows *sqlx.Rows
	if rows, err = r.getDB().QueryxContext(
		ctx, query, filter.Query, filter.Offset, filter.Limit,
	); err != nil {
		return nil, errors.Wrap(err, "error exec query")
	}

	var result = make([]*domain.SearchEntry, 0)
	for rows.Next() {
		entry, err := r.fetch(rows)
		if err != nil {
			return nil, errors.Wrap(err, "error fetch data")
		}

		result = append(result, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error read rows")
	}

	return result, nil
}

func (r *searchRepository) fetch(s scanner) (*domain.SearchEntry, error) {
	var entry domain.SearchEntry
	if err := s.Scan(
		&entry.BoatID,
		&entry.FleetName,
		&entry.ModelName,
		&entry.BuilderName,
	); err != nil {
		return nil, errors.Wrap(err, "can't scan row from query")
	}

	return &entry, nil
}

func (r *searchRepository) getDB() *sqlx.DB {
	return sqlx.NewDb(r.db, "postgres")
}
