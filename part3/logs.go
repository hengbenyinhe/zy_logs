package zy_logs

import (
	"fmt"
	"sync"
	"time"
	"context"
)

//定义日志等级数据类型，便于后面定义日志等级属性
type LogLevel int

var (
	defaultServiceName     =   "default" //定义默认服务名变量
	lm          *LoggerMgr //定义一个LoggerMgr类型的变量
	initOnce    *sync.Once = &sync.Once{} //这个主要是解决多线程调用日志库带来的并发问题
)

//日志记录数据中可能会加入其他字段，例如访问日志会传入用户名等字段
type KeyVal struct {
	key interface{}
	val interface{}
}

type LogField struct {
	kvs []KeyVal
	fieldLock  sync.Mutex //加锁防止并发问题
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

//定义一个日志管理结构
type LoggerMgr struct {
	outputers      []Outputer  //日志输出器
	chanSize       int     //管道缓冲区大小
	level          LogLevel    //日志等级
	logDataChan    chan *LogData  //分配日志管道
	serviceName    string   //产生日志服务名称
	wg             sync.WaitGroup   //阻塞等待日志协程写完才继续执行程序
}

//定义初始化日志方法
func initLogger(level LogLevel, chanSize int, serviceName string) {
	initOnce.Do(func() {
		lm = &LoggerMgr{
			chanSize:chanSize,
			level:level,
			serviceName:serviceName,
			logDataChan:make(chan *LogData,chanSize),
		}
	})
}

//定义对外初始化日志
func InitLogger(level LogLevel,chanSize int,serviceName string)  {
	if chanSize <= 0{
		chanSize = DefaultLogChanSize;
	}

	initLogger(level,chanSize,serviceName)
}

//对外暴露日志函数,将日志分为访问，调试，追踪，普通，警告，错误六种级别
func Access(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf("hi!该函数产生访问级别的日志记录:%v\n",format)
}
func Debug(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf("hi!该函数产生调试级别的日志记录:%v\n",format)
}

func Trace(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf("hi!该函数产生追踪级别的日志记录:%v\n",format)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf("hi!该函数产生普通级别的日志记录:%v\n",format)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf("hi!该函数产生警告级别的日志记录:%v\n",format)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	fmt.Printf("hi!该函数产生错误级别的日志记录:%v\n",format)
}