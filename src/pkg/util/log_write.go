package util

import (
	"fmt"
	"slices"
	"time"
)

var levelComponentColors = []string{
	"\033[38;5;112m", // Hex: #87d700 intlLevel_Info
	"\033[38;5;153m", // Hex: #afd7ff intlLevel_Action
	"\033[38;5;73m",  // Hex: #5fafaf intlLevel_Event
	"\033[38;5;229m", // Hex: #ffffaf intlLevel_Warning
	"\033[38;5;160m", // Hex: #d70000 intlLevel_Error
}

var debugComponentColors = []string{
	"\033[38;5;67m",  // Hex: #5f87af LOG_DEBUG_GENERIC  -> any
	"\033[38;5;208m", // Hex: #ff8700 LOG_DEBUG_TRAFFIC  -> TCP/UDP/HTTP Connections
	"\033[38;5;221m", // Hex: #ffd75f LOG_DEBUG_SERVICE  -> Microservice
	"\033[38;5;191m", // Hex: #d7ff5f LOG_DEBUG_DATABASE -> GORM SQL/CouchDB queries
	"\033[38;5;111m", // Hex: #87afff LOG_DEBUG_API      -> Extra API logging (eg. WebAPI)
	"\033[38;5;62m",  // Hex: #5f5fd7 LOG_DEBUG_ROUTER   -> Spirit Internal
}

const (
	intlLevel_Info int = iota
	intlLevel_Action
	intlLevel_Event
	intlLevel_Warning
	intlLevel_Error
)

const (
	LOG_DEBUG_GENERIC int = iota
	LOG_DEBUG_TRAFFIC
	LOG_DEBUG_SERVICE
	LOG_DEBUG_DATABASE
	LOG_DEBUG_API
	LOG_DEBUG_ROUTER
)

func (l *utilLog) processLog(colorCode string, logType string, prefix string, format string, a ...any) {
	l.writeMutex.Lock()

	formatBuffer := fmt.Sprintf(format, a...)
	finalBuffer := fmt.Sprintf("[%s] [%s] <%s> %s", time.Now().Local().Format("02/01/2006 15:04:05"), logType, prefix, formatBuffer)

	{ // Console Write
		conBuffer := fmt.Sprintf("%s%s\033[0m", colorCode, finalBuffer)
		fmt.Println(conBuffer)
	}

	// File Write
	if l.logFileHandle != nil {
		fileBuffer := fmt.Sprintf("%s\n", finalBuffer)
		l.logFileHandle.WriteString(fileBuffer)
	} else {
		l.unwrittenLogs = append(l.unwrittenLogs, (finalBuffer + "\n"))
	}

	l.writeMutex.Unlock()
}

func (l *utilLog) reconsileLogs() {
	l.writeMutex.Lock()

	for _, logMsg := range l.unwrittenLogs {
		l.logFileHandle.WriteString(logMsg)
	}
	l.unwrittenLogs = slices.Delete(l.unwrittenLogs, 0, len(l.unwrittenLogs))

	l.writeMutex.Unlock()
}

func (l utilLog) Info(prefix string, format string, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Info], "Info", prefix, format, a...)
}

func (l utilLog) Action(prefix string, format string, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Action], "Action", prefix, format, a...)
}

func (l utilLog) Event(prefix string, format string, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Event], "Event", prefix, format, a...)
}

func (l utilLog) Warning(prefix string, format string, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Warning], "Warning", prefix, format, a...)
}

func (l utilLog) Error(prefix string, format string, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Error], "Error", prefix, format, a...)
}

func (l utilLog) Panic(prefix string, format string, err error, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Error], "Fatal", fmt.Sprintf("%s %s", format, err.Error()), format, a...)
	panic(err)
}

func (l utilLog) Debug(component int, prefix string, format string, a ...any) {
	if *Settings.ShowServerDebug {
		l.processLog(debugComponentColors[component], "Debug", prefix, format, a...)
	}
}

func (l utilLog) Trace(component int, prefix string, format string, a ...any) {
	if *Settings.ShowServerDebug {
		l.processLog(debugComponentColors[component], "Trace", prefix, format, a...)
	}
}

func (l *utilLog) Raw(colorCode, format string, a ...any) {
	l.writeMutex.Lock()

	formatBuffer := fmt.Sprintf(format, a...)

	{ // Console Write
		conBuffer := fmt.Sprintf("%s%s\033[0m", colorCode, formatBuffer)
		fmt.Println(conBuffer)
	}

	// File Write
	if l.logFileHandle != nil {
		fileBuffer := fmt.Sprintf("%s\n", formatBuffer)
		l.logFileHandle.WriteString(fileBuffer)
	} else {
		l.unwrittenLogs = append(l.unwrittenLogs, (formatBuffer + "\n"))
	}

	l.writeMutex.Unlock()
}
