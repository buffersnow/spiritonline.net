package db

import (
	"database/sql"
	"fmt"
	"time"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/util"
	"github.com/Masterminds/squirrel"
)

func (s *SQL) ExecRaw(query string, args ...any) (sql.Result, error) {
	start := time.Now()
	res, err := s.e.Exec(query, args...)

	elapsed := time.Since(start)
	var rows int64 = -1
	if err == nil {
		rows, _ = res.RowsAffected()
		s.f.Trace(log.DEBUG_DATABASE, "Exec", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	} else {
		s.f.Warning("Exec", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	}

	return res, err
}

func (s *SQL) Insert(ib squirrel.InsertBuilder) (sql.Result, error) {
	start := time.Now()
	query, args, err := ib.ToSql()
	if err != nil {
		s.f.Error("Insert", "Failed to build SQL query: %v", err)
		return nil, fmt.Errorf("db: sqlx: squirrel: %w", err)
	}

	res, err := s.e.Exec(query, args...)

	elapsed := time.Since(start)
	var rows int64 = -1
	if err == nil {
		rows, _ = res.RowsAffected()
		s.f.Trace(log.DEBUG_DATABASE, "Insert", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	} else {
		s.f.Warning("Insert", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	}

	return res, err
}

func (s *SQL) Update(ub squirrel.UpdateBuilder) (sql.Result, error) {
	start := time.Now()
	query, args, err := ub.ToSql()
	if err != nil {
		s.f.Error("Update", "Failed to build SQL query: %v", err)
		return nil, fmt.Errorf("db: sqlx: squirrel: %w", err)
	}

	res, err := s.e.Exec(query, args...)
	elapsed := time.Since(start)
	var rows int64 = -1
	if err == nil {
		rows, _ = res.RowsAffected()
		s.f.Trace(log.DEBUG_DATABASE, "Update", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	} else {
		s.f.Warning("Update", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	}

	return res, err
}

func (s *SQL) Delete(db squirrel.DeleteBuilder) (sql.Result, error) {
	start := time.Now()
	query, args, err := db.ToSql()
	if err != nil {
		s.f.Error("Delete", "Failed to build SQL query: %v", err)
		return nil, fmt.Errorf("db: sqlx: squirrel: %w", err)
	}

	res, err := s.e.Exec(query, args...)
	elapsed := time.Since(start)
	var rows int64 = -1
	if err == nil {
		rows, _ = res.RowsAffected()
		s.f.Trace(log.DEBUG_DATABASE, "Delete", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	} else {
		s.f.Warning("Delete", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	}

	return res, err
}

func (s *SQL) Get(dest any, sb squirrel.SelectBuilder) error {
	start := time.Now()
	query, args, err := sb.ToSql()
	if err != nil {
		s.f.Error("Get", "Failed to build SQL query: %v", err)
		return fmt.Errorf("db: sqlx: squirrel: %w", err)
	}

	err = s.e.Get(dest, query, args...)
	elapsed := time.Since(start)
	rows := util.CountSQLRows(dest)
	if err == sql.ErrNoRows {
		rows = 0
	}

	if err != nil {
		s.f.Warning("Get", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	} else {
		s.f.Trace(log.DEBUG_DATABASE, "Get", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	}

	return err
}

func (s *SQL) Select(dest any, sb squirrel.SelectBuilder) error {
	start := time.Now()
	query, args, err := sb.ToSql()
	if err != nil {
		s.f.Error("Select", "Failed to build SQL query: %v", err)
		return fmt.Errorf("db: sqlx: squirrel: %w", err)
	}

	err = s.e.Select(dest, query, args...)
	elapsed := time.Since(start)
	rows := util.CountSQLRows(dest)
	if err != nil {
		rows = 0
		s.f.Warning("Select", "<Time: %s> <Rows: None> %v (%s)", elapsed, err, util.FormatSQL(query, args))
	} else {
		s.f.Trace(log.DEBUG_DATABASE, "Select", "<Time: %s> <Rows: %d> %s", elapsed, rows, util.FormatSQL(query, args))
	}

	return err
}
