package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Alekseyt9/ypmetrics/internal/common"
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

func (store *DBStorage) GetCounters(ctx context.Context) (res []common.CounterItem, err error) {
	res = []common.CounterItem{}
	var rows *sql.Rows
	rows, err = store.conn.QueryContext(ctx, "select name, value from counters")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := common.CounterItem{}
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

func (store *DBStorage) SetCounters(ctx context.Context, items []common.CounterItem) error {
	tx, err := store.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		insert into counters(name, value)
		values ($1, $2)
		on conflict (name)
		do update set value = counters.value + EXCLUDED.value
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, item := range items {
		_, err = stmt.ExecContext(ctx, item.Name, item.Value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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

func (store *DBStorage) GetGauges(ctx context.Context) (res []common.GaugeItem, err error) {
	res = []common.GaugeItem{}
	var rows *sql.Rows
	rows, err = store.conn.QueryContext(ctx, "select name, value from gauges")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := common.GaugeItem{}
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

func (store *DBStorage) SetGauges(ctx context.Context, items []common.GaugeItem) error {
	tx, err := store.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		insert into gauges(name, value)
		values ($1, $2)
		on conflict (name)
		do update set value = EXCLUDED.value
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, item := range items {
		_, err = stmt.ExecContext(ctx, item.Name, item.Value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
