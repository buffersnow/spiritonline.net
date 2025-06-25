package lifecycle

import (
	"os"
	"os/signal"
	"syscall"

	"buffersnow.com/spiritonline/pkg/log"
)

/*The only reason this package exists is because round-dependency nonsense*/

func AwaitInterrupt(log *log.Logger) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	if sig := <-c; sig != nil {
		log.Info("Lifecycle", "Captured %v! Stopping Server...", sig)
		os.Exit(0)
	}

	return nil
}
