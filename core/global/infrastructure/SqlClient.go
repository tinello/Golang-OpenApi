package global_infra

import (
	"database/sql"
)

func NewSqlClient(db *sql.DB) *SqlClient {
	return &SqlClient{
		db: db,
	}
}

type SqlClient struct {
	db *sql.DB
}

func (c *SqlClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

func (c *SqlClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

func (c *SqlClient) ExecInTransaction(block func(SqlDatabase) error) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	err = block(newSqlTransaction(tx))
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

type SqlDatabase interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

func newSqlTransaction(tx *sql.Tx) *sqlTransaction {
	return &sqlTransaction{
		tx: tx,
	}
}

type sqlTransaction struct {
	tx *sql.Tx
}

func (t *sqlTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}
