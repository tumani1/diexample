package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/tumani1/diexample/di/sarulabsdi/internal/domain"
)

type autoCompleteRepository struct {
	db *sql.DB
}

func NewAutoCompleteRepository(db *sql.DB) domain.IAutoCompleteRepository {
	return &autoCompleteRepository{
		db: db,
	}
}

func (r *autoCompleteRepository) Find(
	ctx context.Context, filter *domain.AutoCompleteFilter,
) (_ []*domain.AutoCompleteEntry, err error) {
	const query = `
WITH autocomplete_query AS (
	SELECT ym.name "name", 'yacht_models' object_type, similarity(ym.name, $1) score
    FROM yacht_models ym
    WHERE ym.name ILIKE $1 || '%'

    UNION

    SELECT yb.name "name", 'yacht_builders' object_type, similarity(yb.name, $1) score
    FROM yacht_builders yb
    WHERE yb.name ILIKE $1 || '%'
)

SELECT t.name, t.object_type
FROM autocomplete_query t
ORDER BY t.score DESC 
OFFSET $2 LIMIT $3`

	var rows *sqlx.Rows
	if rows, err = r.getDB().QueryxContext(
		ctx, query, filter.Query, filter.Offset, filter.Limit,
	); err != nil {
		return nil, errors.Wrap(err, "error exec query")
	}

	var result = make([]*domain.AutoCompleteEntry, 0)
	for rows.Next() {
		var entry domain.AutoCompleteEntry
		if err := rows.Scan(
			&entry.Name,
			&entry.ObjectType,
		); err != nil {
			return nil, errors.Wrap(err, "can't scan row from query")
		}

		result = append(result, &entry)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error read rows")
	}

	return result, nil
}

func (r *autoCompleteRepository) getDB() *sqlx.DB {
	return sqlx.NewDb(r.db, "postgres")
}
