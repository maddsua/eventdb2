package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/maddsua/eventdb2/storage/model"
	"github.com/maddsua/eventdb2/storage/sqlite/encoders"
	"github.com/maddsua/eventdb2/storage/sqlite/generated"
	_ "github.com/mattn/go-sqlite3"
)

func New(basedir string) (*sqlite, error) {

	dbUrl := url.URL{
		Path: "evdb2.main.db3",
		RawQuery: url.Values{
			"_fk":      []string{"true"},
			"_journal": []string{"WAL"},
		}.Encode(),
	}

	if basedir = strings.ReplaceAll(strings.TrimSpace(basedir), "\\", "/"); basedir != "" && basedir != "/" {

		if _, err := os.Stat(basedir); os.IsNotExist(err) {
			err := os.MkdirAll(basedir, fs.ModeDir)
			if err != nil {
				return nil, fmt.Errorf("unable to create data directory: %v", err)
			}
		}

		dbUrl.Path = filepath.Join(basedir, dbUrl.Path)
	}

	db, err := sql.Open("sqlite3", dbUrl.String())
	if err != nil {
		return nil, fmt.Errorf("unable to open sqlite storage: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("sqlite storage unreachable: %v", err)
	}

	return &sqlite{db: db}, nil
}

type sqlite struct {
	db *sql.DB
}

func (this *sqlite) ID() string {
	return "sqlite"
}

func (this *sqlite) Check() error {
	return this.db.Ping()
}

func (this *sqlite) Close() error {
	return this.db.Close()
}

func (this *sqlite) InsertLogBatch(ctx context.Context, entries []model.LogEntry) error {

	if len(entries) == 0 {
		return errors.New("empty batch")
	}

	tx, err := this.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql.BeginTx: %v", err)
	}

	defer tx.Rollback()

	dbq := generated.New(tx)

	var insert = func(entry model.LogEntry) error {

		metaBuffer, err := encoders.Metadata(entry.Meta).MarshalBinary()
		if err != nil {
			return fmt.Errorf("encoders.Metadata.MarshalBinary: %v", err)
		}

		if err := dbq.InsertLogEntry(ctx, generated.InsertLogEntryParams{
			StreamID: entry.StreamID[:],
			Date:     entry.Date.UnixNano(),
			Level:    string(entry.Level),
			Message:  entry.Message,
			Meta:     metaBuffer,
		}); err != nil {
			return fmt.Errorf("sqlc.InsertLogEntry: %v", err)
		}

		return nil
	}

	for _, entry := range entries {
		if err := insert(entry); err != nil {
			return fmt.Errorf("insert: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("sql.Commit: %v", err)
	}

	return nil
}
