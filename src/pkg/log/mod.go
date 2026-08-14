package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/version"
	"github.com/luxploit/red"
)

type Logger struct {
	mu         *sync.Mutex
	fileName   string
	filePath   string
	fileHandle *os.File
	unwritten  []string
	settings   *LoggerSettings
}

type LoggerVerbosity int

const (
	LoggerVerbosity_Standard LoggerVerbosity = iota
	LoggerVerbosity_Info
	LoggerVerbosity_Trace
	LoggerVerbosity_Debug
)

type LoggerSettings struct {
	Verbosity   LoggerVerbosity `arg:"-v" default:"0" options:"level,maxlevel=3"`
	LogArchival bool            `arg:"--no-archives" default:"true" options:"negated"`
}

var instance = &Logger{mu: &sync.Mutex{}}

func Options(p *settings.Parser) (*LoggerSettings, error) {
	s := &LoggerSettings{}
	return s, p.RegisterModule("config", "Logging", s)
}

func New(r *red.Context, ver *version.BuildTag, o *LoggerSettings) (*Logger, error) {
	log := instance //& this is only a pointer for convinence
	log.settings = o

	tasks := []func() error{}

	log.fileName = fmt.Sprintf("%s.log", ver.GetService())
	log.filePath = fmt.Sprintf("logs/%s.log", ver.GetService())

	log.ToFile("bFXServer - Start Up")
	log.ToFile("SpiritOnline! Build Tag: %s", ver.GetFullTag())
	log.ToFile("Runtime Options: %+v", os.Args[1:])
	log.ToFile("CI by Build Slave: %s", ver.GetCISlave())
	log.ToFile("Start Time: %v", time.Now())

	tasks = append(tasks, log.createLogsFolder)
	if log.settings.LogArchival {
		tasks = append(tasks, log.archiveLog)
	} else {
		log.Warning("Log Provider", "Logfile archival disabled!")
	}

	tasks = append(tasks, log.openLogFile)
	if err := util.Batch(tasks); err != nil {
		return nil, err
	}

	log.reconsileLogs()

	if log.settings.LogArchival {
		go log.archiveLogJob()
	}

	return log, nil
}

func Global() *Logger {
	return instance
}

func (l *Logger) openLogFile() error {
	file, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("log: %w", err)
	}
	l.fileHandle = file
	return nil
}

func (l *Logger) createLogsFolder() error {
	_, err := os.Stat("logs")
	if os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			return fmt.Errorf("log: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("log: %w", err)
	}
	return nil
}
