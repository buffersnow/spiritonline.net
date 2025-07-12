package settings

import (
	"flag"
	"fmt"
	"os"

	"buffersnow.com/spiritonline/pkg/version"
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
	if err != nil {
		return fmt.Errorf("settings: %w", err)
	}

	settings.Runtime = runtimeOptions{
		DisableDB:   *testNoDb,
		DBMigration: *migrateDB,
		LogArchival: !*noLogArchival,
		EnableDebug: *showServerDebug,
		CertsFolder: *certsFolder,
	}

	return nil
}

func (settings *Options) loadEnv(ver *version.BuildTag) func() error {

	return func() error {
		err := godotenv.Load(fmt.Sprintf(".env.%s", ver.GetService()))
		if err != nil {
			return fmt.Errorf("settings: %w", err)
		}

		var cfg Options
		err = env.Parse(&cfg)
		if err != nil {
			return fmt.Errorf("settings: %w", err)
		}

		settings.MySQL = cfg.MySQL
		settings.Spirit = cfg.Spirit
		settings.Service = cfg.Service

		return nil
	}
}
