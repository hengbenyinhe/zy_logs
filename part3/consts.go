package zy_logs

const (
	//日志级别
	LogLevelDebug  LogLevel = iota
	LogLevelTrace
	LogLevelAccess
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

const (
	DefaultLogChanSize = 20000
	SpaceSep           = " "
	ColonSep           = ":"
	LineSep            = "\n"
)
