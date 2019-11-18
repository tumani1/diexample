package domain

import (
	"context"
)

type AutoCompleteEntry struct {
	Name       string `db:"name"`
	ObjectType string `db:"object_type"`
}

type AutoCompleteFilter struct {
	Query  string
	Limit  int64
	Offset int64
}

//go:generate mockgen -source $GOFILE -package mocks -destination mocks/auto_complete.go
type IAutoCompleteRepository interface {
	Find(ctx context.Context, filter *AutoCompleteFilter) ([]*AutoCompleteEntry, error)
}
