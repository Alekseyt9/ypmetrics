package storage

import (
	"context"
	"database/sql"
	"errors"
)

type DBStorage struct {
	conn *sql.DB
}

func NewDBStorage(conn *sql.DB) *DBStorage {
	return &DBStorage{
		conn: conn,
	}
}

func (store *DBStorage) GetCounter(ctx context.Context, name string) (value int64, err error) {
	row := store.conn.QueryRowContext(ctx, `SELECT value FROM counters WHERE name = $1`, name)
	err = row.Scan(&value)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}

	return
}

func (store *DBStorage) SetCounter(ctx context.Context, name string, value int64) (err error) {
	_, err = store.conn.ExecContext(ctx, `
		insert into counters(name, value)
		values ($1, $2)
		on conflict (name)
		do update set value = counters.value + EXCLUDED.value
	`, name, value)
	return
}

func (store *DBStorage) GetCounterAll(ctx context.Context) (res []NameValueCounter, err error) {
	res = []NameValueCounter{}
	var rows *sql.Rows
	rows, err = store.conn.QueryContext(ctx, "select name, value from counters")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := NameValueCounter{}
		if err = rows.Scan(&r.Name, &r.Value); err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (store *DBStorage) GetGauge(ctx context.Context, name string) (value float64, err error) {
	row := store.conn.QueryRowContext(ctx, `SELECT value FROM gauges WHERE name = $1`, name)
	err = row.Scan(&value)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}

	return
}

func (store *DBStorage) SetGauge(ctx context.Context, name string, value float64) (err error) {
	_, err = store.conn.ExecContext(ctx, `
		insert into gauges(name, value)
		values ($1, $2)
		on conflict (name)
		do update set value = EXCLUDED.value
	`, name, value)
	return
}

func (store *DBStorage) GetGaugeAll(ctx context.Context) (res []NameValueGauge, err error) {
	res = []NameValueGauge{}
	var rows *sql.Rows
	rows, err = store.conn.QueryContext(ctx, "select name, value from gauges")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := NameValueGauge{}
		if err = rows.Scan(&r.Name, &r.Value); err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (store *DBStorage) Bootstrap(ctx context.Context) error {
	tx, err := store.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	tx.ExecContext(ctx, `
		CREATE TABLE gauges (
			name varchar(128) PRIMARY KEY,
			value double precision
		);

		CREATE TABLE counters (
			name varchar(128) PRIMARY KEY,
			value bigint
		);
	`) // Индекс не нужет, тк в pg для pk индекс создается автоматически.

	return tx.Commit()
}
