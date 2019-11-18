package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"gitlab.com/igor.tumanov1/theboatscom/internal/domain"
)

type boatAvailableEntry struct {
	BoatID    int64
	Available bool
}

type calendarRepository struct {
	db *sql.DB
}

func NewCalendarRepository(db *sql.DB) domain.ICalendarRepository {
	return &calendarRepository{
		db: db,
	}
}

func (r *calendarRepository) GetAvailabilityByIDs(ctx context.Context, ids []int64) (_ map[int64]bool, err error) {
	const query = `
SELECT boat_id, available 
FROM calendar 
WHERE boat_id = ANY($1) AND date_from <= now()::date AND date_to >= now()::date
`

	var rows *sqlx.Rows
	if rows, err = r.getDB().QueryxContext(ctx, query, pq.Array(&ids)); err != nil {
		return nil, errors.Wrap(err, "error exec query")
	}

	var result = make(map[int64]bool)
	for rows.Next() {
		var entry boatAvailableEntry
		if err := rows.Scan(
			&entry.BoatID,
			&entry.Available,
		); err != nil {
			return nil, errors.Wrap(err, "can't scan row from query")
		}

		result[entry.BoatID] = entry.Available
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error read rows")
	}

	return result, nil
}

func (r *calendarRepository) GetUpcomingAvailabilityDatesByIDs(
	ctx context.Context, ids []int64,
) (_ map[int64]*domain.CalendarEntry, err error) {
	const query = `
WITH ranked_calendar AS (
	SELECT *, ROW_NUMBER() OVER(PARTITION BY boat_id ORDER BY date_from ASC) AS row_number
	FROM calendar 
	WHERE boat_id = ANY($1) AND date_from > now() AND available = true  
)

SELECT boat_id, date_from, date_to, available
FROM ranked_calendar
WHERE row_number = 1
`

	var rows *sqlx.Rows
	if rows, err = r.getDB().QueryxContext(ctx, query, pq.Array(&ids)); err != nil {
		return nil, errors.Wrap(err, "error exec query")
	}

	var result = make(map[int64]*domain.CalendarEntry)
	for rows.Next() {
		var entry domain.CalendarEntry
		if err := rows.Scan(
			&entry.BoatID,
			&entry.DateFrom,
			&entry.DateTo,
			&entry.Available,
		); err != nil {
			return nil, errors.Wrap(err, "can't scan row from query")
		}

		result[entry.BoatID] = &entry
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error read rows")
	}

	return result, nil
}

func (r *calendarRepository) getDB() *sqlx.DB {
	return sqlx.NewDb(r.db, "postgres")
}
