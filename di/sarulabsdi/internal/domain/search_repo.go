package domain

import (
	"context"
)

type SearchEntry struct {
	BoatID      int64  `db:"boat_id"`
	ModelName   string `db:"model_name"`
	BuilderName string `db:"builder_name"`
	FleetName   string `db:"fleet_name"`
}

type SearchFilter struct {
	Query  string
	Limit  int64
	Offset int64
}

//go:generate mockgen -source $GOFILE -package mocks -destination mocks/search.go
type ISearchRepository interface {
	Find(ctx context.Context, filter *SearchFilter) ([]*SearchEntry, error)
}
