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

type SQLSettings struct {
	Enabled  bool `arg:"--no-database" default:"true" options:"negated"`
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

func Options(p *settings.Parser) (*SQLSettings, error) {
	s := &SQLSettings{}
	return s, p.RegisterModule("config", "Database", s)
}

func New(log *log.Logger, o *SQLSettings) (*SQL, error) {
	sql := &SQL{}

	if !o.Enabled {
		log.Warning("Database", "Connection to DB disabled for development purposes!")
		return sql, nil
	}

	log.Action("Database", "Connecting to MySQL database")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		o.Username, o.Password, o.Host, o.Port, o.Database,
	)

	engine, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("db: sqlx: %w", err)
	}

	sql.e = engine
	sql.f = log.Factory("SqlX")

	var version string
	if err := sql.Get(&version, squirrel.Select("version()")); err != nil {
		//& this really is just because its unclear otherwise
		log.Error("Database", "Failed to query MySQL version")
		return nil, fmt.Errorf("db: sqlx: %w", err)
	}

	log.Info("Database", "Connected to \"MySQL Server v%s\"", version)

	return sql, nil
}
