package postgres

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Postgres = Tx{}

type Tx struct {
	db pgx.Tx
}

type TxReq func(tx Tx) error

type TxRunner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type PgxErrorLog func(err error)

func ExecTx(ctx context.Context, log PgxErrorLog, run TxRunner, req TxReq) error {
	pgxTx, err := run.Begin(ctx)
	if err != nil {
		log(err)
		return fmt.Errorf("db begin tx error")
	}
	tx := Tx{
		db: pgxTx,
	}
	defer tx.Rollback(ctx)

	err = req(tx)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		log(err)
		return fmt.Errorf("db commit tx error")
	}

	return nil
}

func (p Tx) Stats() *pgxpool.Stat {
	return nil
}

func (p Tx) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.db.Begin(ctx)
}

func (p Tx) Rollback(ctx context.Context) {
	_ = p.db.Rollback(ctx)
}

func (p Tx) Commit(ctx context.Context) error {
	return p.db.Commit(ctx)
}

func (p Tx) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return p.db.Query(ctx, query, args[:]...)
}

func (p Tx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(ctx, query, args[:]...)
	if err != nil {
		return err
	}
	return pgxscan.DefaultAPI.ScanAll(dest, rows)
}

func (p Tx) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return p.db.Exec(ctx, query, args[:]...)
}

func (p Tx) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return p.db.QueryRow(ctx, query, args[:]...)
}

func (p Tx) Close() {
	p.db.Conn().Close(context.Background())
}
