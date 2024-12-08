package database

import (
	"context"
	"database/sql"
	"time"
)

var defaultTimeout = 5 * time.Second

func OpenDB(dns string, maxOpenConns, maxIdleConns int, connMaxIdleTime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", dns)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil

}
