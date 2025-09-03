package db

import (
	"fmt"
	"reflect"
	"strings"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SQL struct {
	e              *gorm.DB
	f              log.LoggingFactory
	allowMigration bool
}

func New(log *log.Logger, opt *settings.Options) (*SQL, error) {
	sql := &SQL{}

	if opt.Development.EnableDev && opt.Development.DisableDB {
		log.Warning("Database", "Connection to DB disabled for development purposes!")
		return sql, nil
	}

	log.Action("Database", "Connecting to MySQL database")

	gormLog := &CustomGormLogger{f: log.Factory("GORM")}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		opt.MySQL.Username, opt.MySQL.Password, opt.MySQL.Host, opt.MySQL.Port, opt.MySQL.Database,
	)

	engine, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormLog})
	if err != nil {
		return nil, fmt.Errorf("db: gorm: %w", err)
	}

	var version string
	if err := engine.Raw("select version()").Scan(&version).Error; err != nil {
		log.Error("Database", "Failed to query MySQL version")
		return nil, fmt.Errorf("db: gorm: %w", err)
	}

	log.Info("Database", "Connected to \"MySQL Server v%s\"", version)

	sql.e = engine
	sql.f = log.Factory("Database")
	sql.allowMigration = opt.Runtime.DBMigration

	return sql, nil
}

func (s *SQL) MigrateObject(objects ...any) error {

	if !s.allowMigration {
		s.f.Info("Migration", "Migrations will not be run")
		return nil
	}

	names := []string{}
	for _, obj := range objects {
		t := reflect.TypeOf(obj)

		if t.Kind() != reflect.Pointer {
			return fmt.Errorf("db: migration object %v is not a pointer", t.Name())
		}

		t = t.Elem()

		if t.Kind() != reflect.Struct {
			return fmt.Errorf("db: migration object %v is not a struct", t.Name())
		}

		names = append(names, t.Name())
	}

	s.f.Action("Migration", "Running migrations for tables: %s", strings.Join(names, ", "))
	err := s.e.AutoMigrate(objects)

	if err != nil {
		return fmt.Errorf("db: gorm: %w", err)
	}

	return nil
}
