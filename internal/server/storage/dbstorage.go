package storage

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"runtime"

	"github.com/Alekseyt9/ypmetrics/internal/common/items"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // needs for init
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DBStorage struct {
	conn *sql.DB
}

func NewDBStorage(connString string) (*DBStorage, error) {
	conn, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	err = bootstrap(connString)
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		conn: conn,
	}, nil
}

func (store *DBStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	row := store.conn.QueryRowContext(ctx, `SELECT value FROM counters WHERE name = $1`, name)
	var value int64
	err := row.Scan(&value)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}

	return value, err
}

func (store *DBStorage) SetCounter(ctx context.Context, name string, value int64) error {
	_, err := store.conn.ExecContext(ctx, `
		insert into counters(name, value)
		values ($1, $2)
		on conflict (name)
		do update set value = counters.value + EXCLUDED.value
	`, name, value)
	return err
}

func (store *DBStorage) GetCounters(ctx context.Context) ([]items.CounterItem, error) {
	var res = []items.CounterItem{}
	var rows *sql.Rows
	rows, err := store.conn.QueryContext(ctx, "select name, value from counters")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := items.CounterItem{}
		if err = rows.Scan(&r.Name, &r.Value); err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (store *DBStorage) SetCounters(ctx context.Context, items []items.CounterItem) error {
	tx, err := store.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback() //nolint:errcheck //defer

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

func (store *DBStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	row := store.conn.QueryRowContext(ctx, `SELECT value FROM gauges WHERE name = $1`, name)
	var value float64
	err := row.Scan(&value)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}

	return value, err
}

func (store *DBStorage) SetGauge(ctx context.Context, name string, value float64) error {
	_, err := store.conn.ExecContext(ctx, `
		insert into gauges(name, value)
		values ($1, $2)
		on conflict (name)
		do update set value = EXCLUDED.value
	`, name, value)
	return err
}

func (store *DBStorage) GetGauges(ctx context.Context) ([]items.GaugeItem, error) {
	var res = []items.GaugeItem{}
	var rows *sql.Rows
	rows, err := store.conn.QueryContext(ctx, "select name, value from gauges")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := items.GaugeItem{}
		if err = rows.Scan(&r.Name, &r.Value); err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (store *DBStorage) SetGauges(ctx context.Context, items []items.GaugeItem) error {
	tx, err := store.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback() //nolint:errcheck //defer

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

func bootstrap(connString string) error {
	mPath := getMigrationPath()
	m, err := migrate.New(mPath, connString)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (store *DBStorage) Ping(ctx context.Context) error {
	return store.conn.PingContext(ctx)
}

func getMigrationPath() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFilePath)
	migrationsPath := filepath.Join(currentDir, "migrations")
	migrationsPath = filepath.ToSlash(migrationsPath)
	return "file://" + migrationsPath
}
