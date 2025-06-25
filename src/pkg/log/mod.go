package log

import (
	"fmt"
	"os"
	"sync"

	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/util"
)

type Logger struct {
	mu         *sync.Mutex
	debug      bool
	fileName   string
	fileHandle *os.File
	unwritten  []string
}

var instance = &Logger{fileName: "console.log", mu: &sync.Mutex{}}

func New(opt *settings.Options) (*Logger, error) {
	log := instance // this is only a pointer for convinence

	tasks := []func() error{}

	log.debug = *opt.ShowServerDebug

	if *opt.NoLogArchival {
		log.Warning("Log Provider", "Logfile archival disabled!")
	} else {
		tasks = append(tasks, log.archiveLog)
	}

	tasks = append(tasks, log.openLogFile)
	err := util.Batch(tasks)

	log.reconsileLogs()

	if !*opt.NoLogArchival {
		go log.archiveLogJob()
	}

	return log, err
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
