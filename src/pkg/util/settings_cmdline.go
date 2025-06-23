package util

import (
	"flag"
)

func (settings *utilSettings) scanFlags() {

	settings.TestNoDb = flag.Bool("nodb", false, "Run without a database (early-stage testing!)")
	settings.MigrateDB = flag.Bool("migrate", false, "Apply SQL migrations to database (run once!)")
	settings.NoLogCompression = flag.Bool("nologcmpr", false, "Disables logfile daily archive and compressions")
	settings.ShowServerDebug = flag.Bool("debug", false, "Show all debug/developer log")
	settings.ConfigFolder = flag.String("config", "configs", "Server configuration folder")
	settings.CertsFolder = flag.String("certs", "certs", "ECDSA public and private key directory")
	settings.Standalone = flag.Bool("standalone", false, "Turns the microservice into a standalone service")
	settings.ReconnectDelay = flag.Int("rcdelay", 5, "Sets the amount of time between reconnection attempts to the router")

	flag.Parse()
}
