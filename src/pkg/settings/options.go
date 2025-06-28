package settings

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func (settings *Options) loadFlags() error {

	testNoDb := flag.Bool("nodb", false, "Run without a database (early-stage testing!)")
	migrateDB := flag.Bool("migrate", false, "Apply SQL migrations to database (run once!)")
	noLogArchival := flag.Bool("noarchive", false, "Disables logfile daily archive and compressions")
	showServerDebug := flag.Bool("debug", false, "Show all debug/developer log")
	certsFolder := flag.String("certs", "certs", "ECDSA public and private key directory")

	err := flag.CommandLine.Parse(os.Args[1:])

	settings.Runtime = runtimeOptions{
		DisableDB:   *testNoDb,
		DBMigration: *migrateDB,
		LogArchival: !*noLogArchival,
		EnableDebug: *showServerDebug,
		CertsFolder: *certsFolder,
	}

	if err != nil {
		return fmt.Errorf("settings: %w", err)
	}

	return nil
}

func (settings *Options) loadEnv() error {

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("settings: %w", err)
	}

	var cfg Options
	err = env.Parse(&cfg)
	settings = &cfg
	if err != nil {
		return fmt.Errorf("settings: %w", err)
	}

	return nil
}
