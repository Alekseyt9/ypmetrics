package storage

import (
	"context"
	"errors"

	"github.com/Alekseyt9/ypmetrics/internal/common"
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
	rc := common.NewRetryControllerStd(func(err error) bool {
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

func (store *DBRetryStorage) GetCounter(ctx context.Context, name string) (value int64, err error) {
	err = doInRetry(func() error {
		value, err = store.s.GetCounter(ctx, name)
		return err
	})
	return value, err
}

func (store *DBRetryStorage) SetCounter(ctx context.Context, name string, value int64) (err error) {
	err = doInRetry(func() error {
		err = store.s.SetCounter(ctx, name, value)
		return err
	})
	return err
}

func (store *DBRetryStorage) GetCounters(ctx context.Context) (res []common.CounterItem, err error) {
	err = doInRetry(func() error {
		res, err = store.s.GetCounters(ctx)
		return err
	})
	return res, err
}

func (store *DBRetryStorage) SetCounters(ctx context.Context, items []common.CounterItem) (err error) {
	err = doInRetry(func() error {
		err = store.s.SetCounters(ctx, items)
		return err
	})
	return err
}

func (store *DBRetryStorage) GetGauge(ctx context.Context, name string) (value float64, err error) {
	err = doInRetry(func() error {
		value, err = store.s.GetGauge(ctx, name)
		return err
	})
	return value, err
}

func (store *DBRetryStorage) SetGauge(ctx context.Context, name string, value float64) (err error) {
	err = doInRetry(func() error {
		err = store.s.SetGauge(ctx, name, value)
		return err
	})
	return err
}

func (store *DBRetryStorage) GetGauges(ctx context.Context) (res []common.GaugeItem, err error) {
	err = doInRetry(func() error {
		res, err = store.s.GetGauges(ctx)
		return err
	})
	return res, err
}

func (store *DBRetryStorage) SetGauges(ctx context.Context, items []common.GaugeItem) (err error) {
	err = doInRetry(func() error {
		err = store.s.SetGauges(ctx, items)
		return err
	})
	return err
}

func (store *DBRetryStorage) Bootstrap(ctx context.Context) (err error) {
	err = doInRetry(func() error {
		err = store.s.Bootstrap(ctx)
		return err
	})
	return err
}
