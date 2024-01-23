package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const (
	_defaultConnectionAttempts = 10
	_defaultConnectionTimeout  = time.Second

	_maxConnections        = int32(100)
	_minConnections        = int32(10)
	_maxConnectionLifeTime = time.Second * 30
	_maxIdleLifeTime       = time.Second * 60
)

type Postgres interface {
	Stats() *pgxpool.Stat
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Close()
	Select(context.Context, interface{}, string, ...interface{}) error
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	TxRunner
}

func Connect(connectionUrl string) (Postgres, error) {
	connectionAttempts := _defaultConnectionAttempts
	var result *pgxpool.Pool
	var err error
	connectionUrl += fmt.Sprintf(" pool_max_conns=%d pool_min_conns=%d pool_max_conn_lifetime=%v pool_max_conn_idle_time=%v",
		_maxConnections, _minConnections, _maxConnectionLifeTime, _maxIdleLifeTime)
	for connectionAttempts > 0 {
		result, err = pgxpool.New(context.Background(), connectionUrl)
		if err == nil {
			break
		}

		connectionAttempts--
		time.Sleep(_defaultConnectionTimeout)
	}

	if result == nil {
		return nil, err
	}

	return &pool{db: result}, nil
}

type pool struct {
	db *pgxpool.Pool
}

func NewPool(db *pgxpool.Pool) Postgres {
	return pool{
		db: db,
	}
}

func (p pool) Stats() *pgxpool.Stat {
	return p.db.Stat()
}

func (p pool) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.db.Begin(ctx)
}

func (p pool) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return p.db.Query(ctx, query, args[:]...)
}

func (p pool) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return p.db.QueryRow(ctx, query, args[:]...)
}

func (p pool) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(ctx, query, args[:]...)
	if err != nil {
		return err
	}
	return pgxscan.DefaultAPI.ScanAll(dest, rows)
}

func (p pool) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return p.db.Exec(ctx, query, args...)
}

func (p pool) Close() {
	p.db.Close()
}
