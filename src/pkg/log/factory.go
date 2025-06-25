package log

import "fmt"

type LoggingFactory struct {
	prefix  string
	postfix string
}

func Factory(prefix string) *LoggingFactory {
	return &LoggingFactory{prefix: prefix}
}

func FactoryWithPostfix(prefix, postfix string) *LoggingFactory {
	return &LoggingFactory{prefix: prefix, postfix: postfix}
}

func (l *LoggingFactory) ChangePrefix(prefix string, a ...any) {
	l.prefix = fmt.Sprintf(prefix, a...) + " "
}

func (l *LoggingFactory) ChangePostfix(postfix string, a ...any) {
	l.postfix = fmt.Sprintf(postfix, a...) + " "
}

func (l LoggingFactory) Info(prefix string, format string, a ...any) {
	instance.Info(l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Action(prefix string, format string, a ...any) {
	instance.Action(l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Event(prefix string, format string, a ...any) {
	instance.Event(l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Warning(prefix string, format string, a ...any) {
	instance.Warning(l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Error(prefix string, format string, a ...any) {
	instance.Error(l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Panic(prefix string, format string, err error, a ...any) {
	instance.Panic(l.prefix+prefix, l.postfix+format, err, a...)
}

func (l LoggingFactory) Debug(component int, prefix string, format string, a ...any) {
	instance.Debug(component, l.prefix+prefix, l.postfix+format, a...)
}

func (l LoggingFactory) Trace(component int, prefix string, format string, a ...any) {
	instance.Trace(component, l.prefix+prefix, l.postfix+format, a...)
}
