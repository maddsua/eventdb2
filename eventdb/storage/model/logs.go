package model

import (
	"context"
	"database/sql"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type LogOperator interface {
	InsertLogBatch(ctx context.Context, entries []LogEntry) error
	QueryLogs(ctx context.Context, filter LogFilter, page TimePagination) ([]LogEntry, error)

	InsertLogStream(ctx context.Context, stream LogStream) (*LogStream, error)
	GetLogStream(ctx context.Context, id uuid.UUID) (*LogStream, error)
	QueryLogStreams(ctx context.Context, fitler LogStreamFilters, page IndexPagination) ([]LogStream, error)
	ClearLogStreamEntries(ctx context.Context, id uuid.UUID, from time.Time, to time.Time) (int, error)
	SetLogStreamName(ctx context.Context, id uuid.UUID, name string) (*LogStream, error)
	SetLogStreamToken(ctx context.Context, id uuid.UUID, token null.String) (*LogStream, error)
	SetLogStreamNetWhitelist(ctx context.Context, id uuid.UUID, netlist []net.IPNet) (*LogStream, error)
}

type LogEntry struct {
	ID       int64
	StreamID uuid.UUID
	Date     time.Time
	Level    LogLevel
	Message  string
	Meta     map[string]string
}

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
	Regex       *string
	Equal       *string
	NotEqual    *string
	Contains    *string
	NotContains *string
	IsEmpty     *string
}

type LogStream struct {
	ID           uuid.UUID
	Created      time.Time
	Updated      time.Time
	Name         string
	Token        *string
	NetWhitelist []net.IPAddr
}

type LogStreamFilters struct {
	ID           uuid.NullUUID
	NameContains *string
}
