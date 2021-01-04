package domain

import (
	"context"
	"time"
)

type CalendarEntry struct {
	BoatID    int64     `db:"boat_id" json:"boat_id"`
	DateFrom  time.Time `db:"date_from" json:"date_from"`
	DateTo    time.Time `db:"date_to" json:"date_to"`
	Available bool      `db:"available" json:"available"`
}

//go:generate mockgen -source $GOFILE -package mocks -destination mocks/calendar.go
type ICalendarRepository interface {
	GetAvailabilityByIDs(ctx context.Context, ids []int64) (map[int64]bool, error)
	GetUpcomingAvailabilityDatesByIDs(ctx context.Context, ids []int64) (map[int64]*CalendarEntry, error)
}
