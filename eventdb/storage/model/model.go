package model

import (
	"context"
	"database/sql"
	"net"
	"time"

	"github.com/google/uuid"
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

type LogOperator interface {
	InsertLogBatch(ctx context.Context, entries []LogEntry) error
	QueryLogs(ctx context.Context, filter LogFilter, page TimePagination) ([]LogEntry, error)

	InsertLogStream(ctx context.Context, stream LogStream) (*LogStream, error)
	GetLogStream(ctx context.Context, id uuid.UUID) (*LogStream, error)
	QueryLogStreams(ctx context.Context, fitler LogStreamFilters, page IndexPagination) ([]LogStream, error)
	ClearLogStreamEntries(ctx context.Context, id uuid.UUID, from time.Time, to time.Time) (int, error)
	SetLogStreamName(ctx context.Context, id uuid.UUID, name string) (*LogStream, error)
	SetLogStreamToken(ctx context.Context, id uuid.UUID, token sql.NullString) (*LogStream, error)
	SetLogStreamNetWhitelist(ctx context.Context, id uuid.UUID, netlist []net.IPNet) (*LogStream, error)
	SetLogStreamPlatform(ctx context.Context, id uuid.UUID, platform sql.NullString) (*LogStream, error)
}

type LogEntry struct {
	ID       int64
	StreamID uuid.UUID
	Date     time.Time
	Level    LogLevel
	Message  string
	Meta     StringMap
}

type StringMap map[string]string

type LogLevel string

const (
	LogLevelError = "error"
	LogLevelWarn  = "warn"
	LogLevelInfo  = "info"
	LogLevelLog   = "log"
	LogLevelDebug = "debug"
)

type LogFilter struct {
	LogLevel sql.Null[LogLevel]
	StreamID uuid.NullUUID
	Labels   []LogLabelFilter
}

type LogLabelFilter struct {
	Key         string
	Equal       sql.NullString
	NotEqual    sql.NullString
	Contains    sql.NullString
	NotContains sql.NullString
	IsEmpty     sql.NullBool
}

type LogStream struct {
	ID           uuid.UUID
	Created      time.Time
	Updated      time.Time
	Name         string
	Token        sql.NullString
	Platform     sql.NullString
	NetWhitelist []net.IPAddr
}

type LogStreamFilters struct {
	ID           uuid.NullUUID
	NameContains sql.NullString
}
