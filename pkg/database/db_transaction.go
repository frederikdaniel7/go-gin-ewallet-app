package database

import (
	"context"
	"database/sql"
	"log"
)

type transactor struct {
	db *sql.DB
}

type txKey struct{}

type Transactor interface {
	WithinTransaction(context.Context, func(context.Context) error) error
}

type Querier interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, q string) (*sql.Stmt, error)
}

func NewTransaction(db *sql.DB) *transactor {
	return &transactor{
		db: db,
	}
}

func PickQuerier(ctx context.Context, db *sql.DB) Querier {
	var querier Querier
	tx := extractTx(ctx)
	if tx != nil {
		querier = tx
		return querier
	} else {
		querier = db
	} 

	return querier
}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (tr *transactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {

	tx, err := tr.db.Begin()
	if err != nil {
		return err
	}

	err = tFunc(injectTx(ctx, tx))
	if err != nil {

		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("rollback transaction: %v", errRollback)
		}
		return err
	}
	if errCommit := tx.Commit(); errCommit != nil {
		log.Printf("commit transaction: %v", errCommit)
	}
	return nil
}
