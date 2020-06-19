package zy_logs

import (
	"fmt"
	"sync"
	"time"
	"context"
	"path"
)

//定义日志等级数据类型，便于后面定义日志等级属性
//type LogLevel int 将其移动到util

var (
	defaultOutputer   = NewConsoleOutputer()
	defaultServiceName     =   "default" //定义默认服务名变量
	lm          *LoggerMgr //定义一个LoggerMgr类型的变量
	initOnce    *sync.Once = &sync.Once{} //这个主要是解决多线程调用日志库带来的并发问题
	initServiceName string
)

//日志记录数据中可能会加入其他字段，例如访问日志会传入用户名等字段
/*type KeyVal struct {
	key interface{}
	val interface{}
}*/

/*type LogField struct {
	kvs []KeyVal
	fieldLock  sync.Mutex //加锁防止并发问题
}*/

//定义一个日志管理结构
type LoggerMgr struct {
	outputers      []Outputer  //日志输出器
	chanSize       int     //管道缓冲区大小
	level          LogLevel    //日志等级
	logDataChan    chan *LogData  //分配日志管道
	serviceName    string   //产生日志服务名称
	wg             sync.WaitGroup   //阻塞等待日志协程写完才继续执行程序
}

type Outputer interface {
	Write(data *LogData)
	Close()
}

//定义日志数据结构
type LogData struct {
	curTime         time.Time   //日志记录的当前时间
	message         string   //日志信息
	timeStr         string   //日志记录当前时间的格式化
	level           LogLevel  //日志级别
	fileName        string   //产生日志的文件名
	lineNo          int     //产生日志的文件行号
	traceId         string  //追踪id便于分布式的聚合
	serviceName     string  //产生日志的服务名，这里可以在初始化设置
	fields          *LogField //日志信息的其他字段，比如访问日志，传入用户名等字段
}

//定义初始化日志方法
func initLogger(level LogLevel, chanSize int, serviceName string) {
	initServiceName = serviceName

	lm = &LoggerMgr{
		chanSize: chanSize,
		level:level,
		serviceName:serviceName,
	}

	initOnce.Do(func() {
		/*lm = &LoggerMgr{
			chanSize: chanSize,
			level:level,
			serviceName:serviceName,
			logDataChan:make(chan *LogData,chanSize),
		}*/
		lm.logDataChan = make(chan *LogData,chanSize)
		//开启一个异步的协程将管道的日志数写入到输出器中
		lm.wg.Add(1)
		go lm.run()

	})
}

//定义对外初始化日志
func InitLogger(level LogLevel,chanSize int,serviceName string)  {
	if chanSize <= 0{
		chanSize = DefaultLogChanSize;
	}

	initLogger(level,chanSize,serviceName)
}

//对外暴露日志API接口,将日志分为访问，
func Access(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx,LogLevelAccess,format,args...)
}
func Debug(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx,LogLevelDebug,format,args...)
}

func Trace(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx,LogLevelTrace,format,args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx,LogLevelInfo,format,args...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx,LogLevelWarn,format,args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	writeLog(ctx,LogLevelError,format,args...)
}

/*停止日志写入,释放资源*/
func Stop()  {
	close(lm.logDataChan)
	lm.wg.Wait()

	for _,outputer := range lm.outputers{
		outputer.Close()
	}

	initOnce = &sync.Once{}
	lm = nil
}

/*日志记录器*/
func writeLog(ctx context.Context, level LogLevel, formmat string,args ...interface{}) {
	if lm == nil {
		initLogger(LogLevelDebug,DefaultLogChanSize,defaultServiceName)
	}

	now := time.Now()
	nowStr := now.Format("2019-10-11 17:47:05.999")

	fileName, lineNo := GetLineInfo()
	fileName = path.Base(fileName)
	msg := fmt.Sprintf(formmat,args...)

	logData := &LogData{
		message: msg,
		curTime:now,
		timeStr:nowStr,
		level:level,
		fileName:fileName,
		lineNo:lineNo,
		traceId:GetTraceId(ctx),
		serviceName:lm.serviceName,
	}

	//如果为访问日志则可能需要添加其他字段
	if level == LogLevelAccess {
		fields := getFields(ctx)
		if fields != nil {
			logData.fields = fields
		}
	}

	select {
	case lm.logDataChan <- logData:
	default:
		return
	}
}

/*添加输出器*/
func AddOutputer(outputer Outputer) {
	//为保持初始化服务名称只初始化一遍就可以到处使用
	if initServiceName !=""{
		defaultServiceName = initServiceName
	}
	if lm == nil {
		initLogger(LogLevelDebug,DefaultLogChanSize,defaultServiceName)
	}
	lm.outputers = append(lm.outputers,outputer)
	return
}

/*从管道中获取日志数据，并且打印到相应的输出器*/
func (l *LoggerMgr) run() {
	for data := range l.logDataChan {
		if len(l.outputers) == 0 {
			defaultOutputer.Write(data)
			continue //跳出这次循环执行下一次循环
		}
		for _, outputer := range l.outputers {
			outputer.Write(data)
		}

	}

	l.wg.Done()
}