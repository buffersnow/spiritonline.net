package log

import (
	"fmt"
	"slices"
	"time"
)

func (l *Logger) processLog(colorCode string, logType string, prefix string, format string, a ...any) {
	l.mu.Lock()

	formatBuffer := fmt.Sprintf(format, a...)
	finalBuffer := fmt.Sprintf("[%s] [%s] <%s> %s", time.Now().Local().Format("02/01/2006 15:04:05"), logType, prefix, formatBuffer)

	// File Write
	if l.fileHandle != nil {
		fileBuffer := fmt.Sprintf("%s\n", finalBuffer)
		l.fileHandle.WriteString(fileBuffer)
	} else {
		l.unwritten = append(l.unwritten, (finalBuffer + "\n"))
	}

	{ // Console Write
		conBuffer := fmt.Sprintf("%s%s\033[0m", colorCode, finalBuffer)
		fmt.Println(conBuffer)
	}

	l.mu.Unlock()
}

func (l *Logger) reconsileLogs() {
	l.mu.Lock()

	for _, logMsg := range l.unwritten {
		l.fileHandle.WriteString(logMsg)
	}
	l.unwritten = slices.Delete(l.unwritten, 0, len(l.unwritten))

	l.mu.Unlock()
}

func (l *Logger) Info(prefix string, format string, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Info], "Info", prefix, format, a...)
}

func (l *Logger) Action(prefix string, format string, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Action], "Action", prefix, format, a...)
}

func (l *Logger) Event(prefix string, format string, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Event], "Event", prefix, format, a...)
}

func (l *Logger) Warning(prefix string, format string, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Warning], "Warning", prefix, format, a...)
}

func (l *Logger) Error(prefix string, format string, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Error], "Error", prefix, format, a...)
}

func (l *Logger) Panic(prefix string, format string, err error, a ...any) {
	l.processLog(levelComponentColors[intlLevel_Error], "Fatal", fmt.Sprintf("%s %s", format, err.Error()), format, a...)
	panic(err)
}

func (l *Logger) Debug(component int, prefix string, format string, a ...any) {
	if l.debug {
		l.processLog(debugComponentColors[component], "Debug", prefix, format, a...)
	}
}

func (l *Logger) Trace(component int, prefix string, format string, a ...any) {
	l.processLog(debugComponentColors[component], "Trace", prefix, format, a...)
}

func (l *Logger) Raw(colorCode, format string, a ...any) {
	l.mu.Lock()

	formatBuffer := fmt.Sprintf(format, a...)

	{ // Console Write
		conBuffer := fmt.Sprintf("%s%s\033[0m", colorCode, formatBuffer)
		fmt.Println(conBuffer)
	}

	// File Write
	if l.fileHandle != nil {
		fileBuffer := fmt.Sprintf("%s\n", formatBuffer)
		l.fileHandle.WriteString(fileBuffer)
	} else {
		l.unwritten = append(l.unwritten, (formatBuffer + "\n"))
	}

	l.mu.Unlock()
}

func (l *Logger) ToFile(format string, a ...any) {
	l.mu.Lock()

	formatBuffer := fmt.Sprintf(format, a...)

	// File Write
	if l.fileHandle != nil {
		fileBuffer := fmt.Sprintf("%s\n", formatBuffer)
		l.fileHandle.WriteString(fileBuffer)
	} else {
		l.unwritten = append(l.unwritten, (formatBuffer + "\n"))
	}

	l.mu.Unlock()
}
