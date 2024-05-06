package storage

import (
	"context"
	"database/sql"
)

type DbStorage struct {
	conn *sql.DB
}

func NewDbStorage(conn *sql.DB) *DbStorage {
	return &DbStorage{
		conn: conn,
	}
}

func (store *DbStorage) GetCounter(ctx context.Context, name string) (value int64, err error) {
	row := store.conn.QueryRowContext(ctx, `SELECT value FROM counters WHERE name = $1`, name)
	err = row.Scan(&value)
	return
}

func (store *DbStorage) SetCounter(ctx context.Context, name string, value int64) (err error) {
	_, err = store.conn.ExecContext(ctx, `
		insert into counters(name, value)
		values ($1, $2)
		on conflict (name)
		do update set value = counters.value + EXCLUDED.value
	`, name, value)
	return
}

func (store *DbStorage) GetCounterAll(ctx context.Context) (res []NameValueCounter, err error) {
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
	return res, nil
}

func (store *DbStorage) GetGauge(ctx context.Context, name string) (value float64, err error) {
	row := store.conn.QueryRowContext(ctx, `SELECT value FROM gauges WHERE name = $1`, name)
	err = row.Scan(&value)
	return
}

func (store *DbStorage) SetGauge(ctx context.Context, name string, value float64) (err error) {
	_, err = store.conn.ExecContext(ctx, `
		insert into gauges(name, value)
		values ($1, $2)
		on conflict (name)
		do update set value = EXCLUDED.value
	`, name, value)
	return
}

func (store *DbStorage) GetGaugeAll(ctx context.Context) (res []NameValueGauge, err error) {
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
	return res, nil
}

func (store *DbStorage) Bootstrap(ctx context.Context) error {
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
