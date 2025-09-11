package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/util"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type SQLTransaction struct {
	c context.Context
	e *sqlx.Tx
	f *log.LoggingFactory
}

func (s *SQL) Begin() (*SQLTransaction, error) {
	ctx := context.Background()

	tx, err := s.e.BeginTxx(ctx, nil)
	if err != nil {
		s.f.Error("Transaction", "Failed to begin transaction: %v", err)
		return nil, err
	}

	s.f.Trace(log.DEBUG_DATABASE, "Transaction", "Started transaction")

	return &SQLTransaction{
		e: tx,
		f: &s.f,
		c: ctx,
	}, nil
}

func (s *SQLTransaction) Commit() error {
	err := s.e.Commit()
	if err != nil {
		s.f.Error("Transaction", "Failed to commit transaction: %v", err)
		return err
	}

	s.f.Trace(log.DEBUG_DATABASE, "Transaction", "Committing transaction")
	return nil
}

func (s *SQLTransaction) Rollback() error {

	err := s.e.Rollback()
	if err != nil {
		s.f.Error("Transaction", "Failed to rollback transaction: %v", err)
		return err
	}

	s.f.Trace(log.DEBUG_DATABASE, "Transaction", "Rolling back transaction")
	return nil
}

func (s *SQLTransaction) ExecRaw(query string, args ...any) (sql.Result, error) {
	start := time.Now()
	res, err := s.e.ExecContext(s.c, query, args...)

	elapsed := time.Since(start)
	var rows int64 = -1
	if err == nil {
		rows, _ = res.RowsAffected()
		s.f.Trace(log.DEBUG_DATABASE, "Exec", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	} else {
		s.f.Error("Exec", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	}

	return res, err
}

func (s *SQLTransaction) Insert(ib squirrel.InsertBuilder) (sql.Result, error) {
	start := time.Now()
	query, args, err := ib.ToSql()
	if err != nil {
		s.f.Error("Insert", "Failed to build SQL query: %v", err)
		return nil, fmt.Errorf("db: sqlx: squirrel: %w", err)
	}

	res, err := s.e.ExecContext(s.c, query, args...)

	elapsed := time.Since(start)
	var rows int64 = -1
	if err == nil {
		rows, _ = res.RowsAffected()
		s.f.Trace(log.DEBUG_DATABASE, "Insert", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	} else {
		s.f.Error("Insert", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	}

	return res, err
}

func (s *SQLTransaction) Update(ub squirrel.UpdateBuilder) (sql.Result, error) {
	start := time.Now()
	query, args, err := ub.ToSql()
	if err != nil {
		s.f.Error("Update", "Failed to build SQL query: %v", err)
		return nil, fmt.Errorf("db: sqlx: squirrel: %w", err)
	}

	res, err := s.e.ExecContext(s.c, query, args...)
	elapsed := time.Since(start)
	var rows int64 = -1
	if err == nil {
		rows, _ = res.RowsAffected()
		s.f.Trace(log.DEBUG_DATABASE, "Update", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	} else {
		s.f.Error("Update", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	}

	return res, err
}

func (s *SQLTransaction) Delete(db squirrel.DeleteBuilder) (sql.Result, error) {
	start := time.Now()
	query, args, err := db.ToSql()
	if err != nil {
		s.f.Error("Delete", "Failed to build SQL query: %v", err)
		return nil, fmt.Errorf("db: sqlx: squirrel: %w", err)
	}

	res, err := s.e.ExecContext(s.c, query, args...)
	elapsed := time.Since(start)
	var rows int64 = -1
	if err == nil {
		rows, _ = res.RowsAffected()
		s.f.Trace(log.DEBUG_DATABASE, "Delete", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	} else {
		s.f.Error("Delete", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	}

	return res, err
}

func (s *SQLTransaction) Get(dest any, sb squirrel.SelectBuilder) error {
	start := time.Now()
	query, args, err := sb.ToSql()
	if err != nil {
		s.f.Error("Get", "Failed to build SQL query: %v", err)
		return fmt.Errorf("db: sqlx: squirrel: %w", err)
	}

	err = s.e.GetContext(s.c, dest, query, args...)
	elapsed := time.Since(start)
	rows := int64(1)
	if err == sql.ErrNoRows {
		rows = 0
	}

	if err != nil {
		s.f.Error("Get", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	} else {
		s.f.Trace(log.DEBUG_DATABASE, "Get", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	}

	return err
}

func (s *SQLTransaction) Select(dest any, sb squirrel.SelectBuilder) error {
	start := time.Now()
	query, args, err := sb.ToSql()
	if err != nil {
		s.f.Error("Select", "Failed to build SQL query: %v", err)
		return fmt.Errorf("db: sqlx: squirrel: %w", err)
	}

	err = s.e.SelectContext(s.c, dest, query, args...)
	elapsed := time.Since(start)
	rows := util.CountSQLRows(dest)
	if err != nil {
		rows = 0
		s.f.Error("Select", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	} else {
		s.f.Trace(log.DEBUG_DATABASE, "Select", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	}

	return err
}
