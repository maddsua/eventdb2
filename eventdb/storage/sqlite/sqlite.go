package sqlite

import (
	"database/sql"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
