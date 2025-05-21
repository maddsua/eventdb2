package model

import "github.com/guregu/null"

type Storage interface {
	ID() string

	LogOperator

	Check() error
	Close() error
}

type TimePagination struct {
	Before null.Time
	After  null.Time
}

type IndexPagination struct {
	Limit  null.Int
	Offser null.Int
}
