package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/version"
)

type Logger struct {
	mu         *sync.Mutex
	debug      bool
	fileName   string
	fileHandle *os.File
	unwritten  []string
}

var instance = &Logger{fileName: "console.log", mu: &sync.Mutex{}}

func New(ver *version.BuildTag, opt *settings.Options) (*Logger, error) {
	log := instance //& this is only a pointer for convinence

	tasks := []func() error{}

	log.debug = opt.Development.EnableDebug

	log.ToFile("SpiritOnline! Build Tag: %s", ver.GetFullTag())
	log.ToFile("Runtime Options: %+v", opt.Runtime)
	log.ToFile("CI by Build Slave: %s", ver.GetCISlave())
	log.ToFile("Start Time: %v", time.Now())

	if opt.Runtime.LogArchival {
		tasks = append(tasks, log.archiveLog)
	} else {
		log.Warning("Log Provider", "Logfile archival disabled!")
	}

	tasks = append(tasks, log.openLogFile)
	if err := util.Batch(tasks); err != nil {
		return nil, err
	}

	log.reconsileLogs()

	if opt.Runtime.LogArchival {
		go log.archiveLogJob()
	}

	return log, nil
}

func Global() *Logger {
	return instance
}

func (l *Logger) openLogFile() error {
	file, err := os.OpenFile(l.fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("log: %w", err)
	}
	l.fileHandle = file
	return nil
}
