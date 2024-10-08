package storage

import (
	"context"
	"errors"

	"github.com/Alekseyt9/ypmetrics/internal/common/items"
	"github.com/Alekseyt9/ypmetrics/pkg/retry"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBRetryStorage struct {
	s *DBStorage
}

func NewDBRetryStorage(s *DBStorage) *DBRetryStorage {
	return &DBRetryStorage{s: s}
}

func doInRetry(f func() error) error {
	rc := retry.NewControllerStd(func(err error) bool {
		code, ok := extractErrorCode(err)
		if !ok {
			return false
		}
		return pgerrcode.IsConnectionException(code)
	})
	return rc.Do(f)
}

func extractErrorCode(err error) (string, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code, true
	}
	return "", false
}

func (store *DBRetryStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	var value int64
	err := doInRetry(func() error {
		var err error
		value, err = store.s.GetCounter(ctx, name)
		return err
	})
	return value, err
}

func (store *DBRetryStorage) SetCounter(ctx context.Context, name string, value int64) error {
	err := doInRetry(func() error {
		err := store.s.SetCounter(ctx, name, value)
		return err
	})
	return err
}

func (store *DBRetryStorage) GetCounters(ctx context.Context) ([]items.CounterItem, error) {
	var res []items.CounterItem
	err := doInRetry(func() error {
		var err error
		res, err = store.s.GetCounters(ctx)
		return err
	})
	return res, err
}

func (store *DBRetryStorage) SetCounters(ctx context.Context, items []items.CounterItem) error {
	err := doInRetry(func() error {
		err := store.s.SetCounters(ctx, items)
		return err
	})
	return err
}

func (store *DBRetryStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	var value float64
	err := doInRetry(func() error {
		var err error
		value, err = store.s.GetGauge(ctx, name)
		return err
	})
	return value, err
}

func (store *DBRetryStorage) SetGauge(ctx context.Context, name string, value float64) error {
	err := doInRetry(func() error {
		err := store.s.SetGauge(ctx, name, value)
		return err
	})
	return err
}

func (store *DBRetryStorage) GetGauges(ctx context.Context) ([]items.GaugeItem, error) {
	var res []items.GaugeItem
	err := doInRetry(func() error {
		var err error
		res, err = store.s.GetGauges(ctx)
		return err
	})
	return res, err
}

func (store *DBRetryStorage) SetGauges(ctx context.Context, items []items.GaugeItem) error {
	err := doInRetry(func() error {
		err := store.s.SetGauges(ctx, items)
		return err
	})
	return err
}

func (store *DBRetryStorage) Ping(ctx context.Context) error {
	err := doInRetry(func() error {
		err := store.s.Ping(ctx)
		return err
	})
	return err
}
