package zy_logs

import (
	"bytes"
	"fmt"
	"runtime"
)

type LogLevel int
/*获取日志等级字符串*/
func getLevelText(level LogLevel) string{
	switch level {
	case LogLevelAccess:
		return "ACCESS"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	}
	return "UNKNOWN"
}
/*根据日志等级字符串返回日志等级*/
func GetLogLevel(level string) LogLevel {
	switch level {
	case "debug":
		return LogLevelDebug
	case "trace":
		return LogLevelTrace
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	}
	return LogLevelDebug
}

/*获取生成日志的文件名*/
func GetLineInfo() (fileName string,lineNo int) {
	_, fileName, lineNo, _ = runtime.Caller(3)
	return
}
/*将字段写入到buffer缓存进行拼接*/
func writeField(buffer *bytes.Buffer,field,sep string)  {
	buffer.WriteString(field)
	buffer.WriteString(sep)
}

/*将结构体的日志数据转化为字节数组*/
func (l *LogData)Bytes() []byte {
	var buffer bytes.Buffer
	levelStr := getLevelText(l.level)

	writeField(&buffer,l.timeStr,SpaceSep)
	writeField(&buffer,levelStr,SpaceSep)
	writeField(&buffer,l.serviceName,SpaceSep)

	writeField(&buffer,l.fileName,ColonSep)
	writeField(&buffer,fmt.Sprintf("%d",l.lineNo),SpaceSep)
	writeField(&buffer,l.traceId,SpaceSep)
	if l.level == LogLevelAccess && l.fields != nil {
		for _,field := range l.fields.kvs {
			writeField(&buffer, fmt.Sprintf("%v=%v",field.key,field.val),SpaceSep)
		}
	}

	writeField(&buffer,l.message,LineSep)

	return buffer.Bytes()

}
/*根据日志级别获取不同颜色*/
func getLevelColor(level LogLevel) Color {
	switch level {
	case LogLevelAccess:
		return Blue
	case LogLevelDebug:
		return White
	case LogLevelTrace:
		return Cyan
	case LogLevelInfo:
		return Green
	case LogLevelWarn:
		return Yellow
	case LogLevelError:
		return Red
	}
	return Magenta
}

/*日志文件的切分时段*/
type LogFileSeg int
/*获取日志等级字符串*/
func getSegText(seg LogFileSeg) string{
	switch seg {
	case YearSeg:
		return "year"
	case MonthSeg:
		return "month"
	case WeekSeg:
		return "week"
	case DaySeg:
		return "day"
	case HourSeg:
		return "hour"
	}
	return "hour"
}
/*根据日志等级字符串返回日志等级*/
func GetFileSeg(seg string) LogFileSeg {
	switch seg {
	case "year":
		return YearSeg
	case "month":
		return MonthSeg
	case "week":
		return WeekSeg
	case "day":
		return DaySeg
	}
	return HourSeg
}
