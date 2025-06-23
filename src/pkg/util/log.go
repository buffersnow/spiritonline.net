package util

import (
	"fmt"
	"os"
	"sync"
)

type utilLog struct {
	logFileName   string
	logFileHandle *os.File
	writeMutex    *sync.Mutex
	unwrittenLogs []string
}

var Log utilLog = utilLog{logFileName: "console.log", writeMutex: &sync.Mutex{}}

// !HAS! to be a go routine so that we don't deadlock main lol
func (l utilLog) Initialize() {
	go func() {
		settingsWg.Wait()

		if *Settings.NoLogCompression {
			l.Warning("Startup", "Logfile compression disabled!")
		} else {
			l.compressLog()
		}

		l.openLogFile()

		l.reconsileLogs()

		if !*Settings.NoLogCompression {
			go l.compressLogJob()
		}
	}()
}

func (l *utilLog) openLogFile() {
	file, err := os.OpenFile(l.logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create log file %s", err.Error()))
	}
	l.logFileHandle = file
}
