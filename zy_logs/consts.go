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
	DefaultLogChanSize = 20000 //默认通道大小
	SpaceSep           = " " //空格分隔符
	ColonSep           = ":"  //冒号分隔符
	LineSep            = "\n" //换行分隔符
)

const (
	//日志控制台输出颜色，经测试貌似只对linux操作系统有效
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const (
	YearSeg LogFileSeg = iota
	MonthSeg
	WeekSeg
	DaySeg
	HourSeg
)


