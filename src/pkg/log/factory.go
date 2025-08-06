package log

import "fmt"

type LoggingFactory struct {
	log     *Logger
	prefix  string
	postfix string
}

func (l *Logger) Factory(prefix string) LoggingFactory {
	return LoggingFactory{log: l, prefix: prefix + " "}
}

func (l *Logger) FactoryWithPostfix(prefix, postfix string) LoggingFactory {
	return LoggingFactory{log: l, prefix: prefix + " ", postfix: postfix + " "}
}

func (l *LoggingFactory) ChangePrefix(prefix string, a ...any) {
	l.prefix = fmt.Sprintf(prefix, a...) + " "
}

func (l *LoggingFactory) ChangePostfix(postfix string, a ...any) {
	l.postfix = fmt.Sprintf(postfix, a...) + " "
}

func (l LoggingFactory) Info(prefix string, format string, a ...any) {
	l.log.Info(l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Action(prefix string, format string, a ...any) {
	l.log.Action(l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Event(prefix string, format string, a ...any) {
	l.log.Event(l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Warning(prefix string, format string, a ...any) {
	l.log.Warning(l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Error(prefix string, format string, a ...any) {
	l.log.Error(l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Panic(prefix string, format string, err error, a ...any) {
	l.log.Panic(l.prefix+prefix, l.postfix+format, err, a...)
}

func (l LoggingFactory) Debug(component int, prefix string, format string, a ...any) {
	l.log.Debug(component, l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Trace(component int, prefix string, format string, a ...any) {
	l.log.Trace(component, l.prefix+prefix, l.postfix+format, a...)
}
