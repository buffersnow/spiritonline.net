package util

import "fmt"

type LogFactory struct {
	prefix  string
	postfix string
}

func (utilLog) Instance(prefix string, a ...any) *LogFactory {
	return &LogFactory{prefix: fmt.Sprintf(prefix, a...) + " "}
}

func (l *LogFactory) ChangePrefix(prefix string, a ...any) {
	l.prefix = fmt.Sprintf(prefix, a...) + " "
}

func (l *LogFactory) ChangePostfix(postfix string, a ...any) {
	l.postfix = fmt.Sprintf(postfix, a...) + " "
}

func (l LogFactory) Info(prefix string, format string, a ...any) {
	Log.Info(l.prefix+prefix, l.postfix+format, a...)
}

func (l LogFactory) Action(prefix string, format string, a ...any) {
	Log.Action(l.prefix+prefix, l.postfix+format, a...)
}

func (l LogFactory) Event(prefix string, format string, a ...any) {
	Log.Event(l.prefix+prefix, l.postfix+format, a...)
}

func (l LogFactory) Warning(prefix string, format string, a ...any) {
	Log.Warning(l.prefix+prefix, l.postfix+format, a...)
}

func (l LogFactory) Error(prefix string, format string, a ...any) {
	Log.Error(l.prefix+prefix, l.postfix+format, a...)
}

func (l LogFactory) Panic(prefix string, format string, err error, a ...any) {
	Log.Panic(l.prefix+prefix, l.postfix+format, err, a...)
}

func (l LogFactory) Debug(component int, prefix string, format string, a ...any) {
	Log.Debug(component, l.prefix+prefix, l.postfix+format, a...)
}

func (l LogFactory) Trace(component int, prefix string, format string, a ...any) {
	Log.Trace(component, l.prefix+prefix, l.postfix+format, a...)
}
