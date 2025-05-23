package model

import (
	"time"
)

type Storage interface {
	ID() string

	LogOperator

	Check() error
	Close() error
}

type TimePagination struct {
	FromDate  *time.Time
	UntilDate *time.Time
}

type IndexPagination struct {
	Limit  *int
	Offser *int
}
