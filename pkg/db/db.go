package db

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema string = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT '' CHECK (length(repeat) <= 128)
	);

CREATE INDEX idx_scheduler_date ON scheduler (date);
`

func Init(dbFile string) error {
	var install bool

	_, err := os.Stat(dbFile)
	if errors.Is(err, os.ErrNotExist) {
		install = true
	} else if err != nil {
		return err
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		if _, err = db.Exec(schema); err != nil {
			return err
		}
	}
	return nil
}
