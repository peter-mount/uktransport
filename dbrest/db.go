package dbrest

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type DB struct {
	PostgresUri string `yaml:"url"`
	MaxOpen     int    `yaml:"maxOpen"`
	MaxIdle     int    `yaml:"maxIdle"`
	MaxLifetime int    `yaml:"maxLifetime"`
	db          *sql.DB
}

func (d *DB) start() error {
	if d.db != nil {
		return nil
	}

	db, err := sql.Open("postgres", d.PostgresUri)
	if err != nil {
		return err
	}
	d.db = db

	if d.MaxOpen < 0 {
		d.MaxOpen = 1
	}

	if d.MaxIdle < 0 {
		d.MaxIdle = 1
	} else if d.MaxIdle > d.MaxOpen {
		d.MaxIdle = d.MaxOpen
	}

	d.db.SetMaxOpenConns(d.MaxOpen)
	d.db.SetMaxIdleConns(d.MaxIdle)

	if d.MaxLifetime > 0 {
		d.db.SetConnMaxLifetime(time.Second * time.Duration(d.MaxLifetime))
	}

	return nil
}

func (d *DB) stop() {
	if d.db != nil {
		_ = d.db.Close()
		d.db = nil
	}
}

func (s *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	r, e := s.db.Exec(query, args...)
	return r, e
}

func (s *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	r, e := s.db.Query(query, args...)
	return r, e
}

func (s *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

func (s *DB) Prepare(sql string) (*sql.Stmt, error) {
	return s.db.Prepare(sql)
}

func (s *DB) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}
