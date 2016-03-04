package api

type logger interface {
	Debug(fmt string, v ...interface{})
	Info(fmt string, v ...interface{})
	Warning(fmt string, v ...interface{})
	Error(fmt string, v ...interface{})
}

var objectLogger logger = nil

func SetLogger(l logger) {
	objectLogger = l
}

func Debug(fmt string, v ...interface{}) {
	if objectLogger != nil {
		objectLogger.Debug(fmt, v)
	}
}

func Info(fmt string, v ...interface{}) {
	if objectLogger != nil {
		objectLogger.Info(fmt, v)
	}
}

func Warning(fmt string, v ...interface{}) {
	if objectLogger != nil {
		objectLogger.Warning(fmt, v)
	}
}

func Error(fmt string, v ...interface{}) {
	if objectLogger != nil {
		objectLogger.Info(fmt, v)
	}
}
