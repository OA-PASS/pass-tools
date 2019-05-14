package log

import (
	logger "log"
	"os"
)

const (
	INFO = iota
	DEBUG
	TRACE
)

type Printer interface {
	Printf(format string, a ...interface{})
}

func New(level int) Instance {

	l := Instance{
		Warn: logger.New(os.Stdout, "WARN: ", logger.Flags()),
	}

	if level >= INFO {
		l.Info = logger.New(os.Stdout, "INFO: ", logger.Flags())
	}

	if level >= DEBUG {
		l.Debug = logger.New(os.Stdout, "DEBUG: ", logger.Flags())
	}

	if level >= TRACE {
		l.Trace = logger.New(os.Stdout, "TRACE: ", logger.Flags())
	}

	return l
}

type Instance struct {
	Warn  Printer
	Info  Printer
	Debug Printer
	Trace Printer
}

func (l Instance) Printf(format string, a ...interface{}) {
	if l.Info == nil {
		logger.Printf(format, a...)
		return
	}
	l.Info.Printf(format, a...)
}

func (l Instance) Warnf(format string, a ...interface{}) {
	if l.Warn == nil {
		l.Printf(format, a...)
	}
	l.Warn.Printf(format, a ...)
}

func (l Instance) Debugf(format string, a ...interface{}) {
	if l.Debug != nil {
		l.Debug.Printf(format, a...)
	}
}

func (l Instance) Tracef(format string, a ...interface{}) {
	if l.Trace != nil {
		l.Trace.Printf(format, a...)
	}
}



