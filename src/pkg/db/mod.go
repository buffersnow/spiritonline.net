package db

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

type SQL struct {
	e *sqlx.DB
	f log.LoggingFactory
}

func New(log *log.Logger, opt *settings.Options) (*SQL, error) {
	sql := &SQL{}

	if opt.Development.EnableDev && opt.Development.DisableDB {
		log.Warning("Database", "Connection to DB disabled for development purposes!")
		return sql, nil
	}

	log.Action("Database", "Connecting to MySQL database")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		opt.MySQL.Username, opt.MySQL.Password, opt.MySQL.Host, opt.MySQL.Port, opt.MySQL.Database,
	)

	engine, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("db: sqlx: %w", err)
	}

	var version string
	if err := sql.Get(&version, squirrel.Select("version()")); err != nil {
		//& this really is just because its unclear otherwise
		log.Error("Database", "Failed to query MySQL version")
		return nil, fmt.Errorf("db: sqlx: %w", err)
	}

	log.Info("Database", "Connected to \"MySQL Server v%s\"", version)

	sql.e = engine
	sql.f = log.Factory("SqlX")

	return sql, nil
}
