package model

type Storage interface {
	ID() string

	LogOperator

	Check() error
	Close() error
}
