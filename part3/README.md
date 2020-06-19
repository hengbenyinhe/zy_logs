# 原型实现-实现日志公共方法

在[上一个章节][第二章]已经定义了一些要用到的数据结构，如果之后的开发过程中可能还需要定义数据结构，到时候会具体讲解的

这章节就实现日志库需要用到的函数或者方法，首先可以肯定的是对外暴露的日志级别函数。

## 定义对外暴露的日志级别函数

在logs.go文件中加入如下代码
```go
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
```

我们测试一下，这些函数是否正常访问。创建一个logs_test.go文件，这里对测试的知识点不做过多的讲解，如果单元测试知识没掌握的朋友可以去查看相关的教程

```go
    package zy_logs
    
    import (
    	"context"
    	"testing"
    )
    
    func TestLoggerFunc(t *testing.T) {
    	Access(context.Background(), "Access")
    	Debug(context.Background(), "Debug")
    	Trace(context.Background(), "Trace")
    	Info(context.Background(), "Info")
    	Warn(context.Background(), "Warn")
    	Error(context.Background(), "Error")
    }
```

然后在终端中执行go test命令,可以看到输出结果如下：
```text
    hi!该函数产生访问级别的日志记录:Access
    hi!该函数产生调试级别的日志记录:Debug
    hi!该函数产生追踪级别的日志记录:Trace
    hi!该函数产生普通级别的日志记录:Info
    hi!该函数产生警告级别的日志记录:Warn
    hi!该函数产生错误级别的日志记录:Error
    PASS
    ok      zy_logs 0.002s
```

![alt 测试结果]("./docImage/testRes.png")
<img src="./docImage/testRes.png" style="zoom:70%" />

则说明定义的函数可以正常使用，之后每当完成一个功能的单元测试就不再此教程中演示了，因为大同小异大家可以以此类推。

接下来开始初始化日志库，在consts.go 文件中定义日志数据结构用到的分隔符常量
```go
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

```
还需要定义一个日志库管理结构体，在logs.go文件中加入一下代码
```go
//定义一个日志管理结构
type LoggerMgr struct {
	outputers      []Outputer  //日志输出器
	chanSize       int     //管道缓冲区大小
	level          LogLevel    //日志等级
	logDataChan    chan *LogData  //分配日志管道
	serviceName    string   //产生日志服务名称
	wg             sync.WaitGroup   //阻塞等待日志协程写完才继续执行程序
}
```
 紧接着定义一个LoggerMgr 类型的变量lm，定义默认服务名变量defaultServiceName
 定义初始化日志方法initLogger, 
```go
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
```
对外公开初始化日志方法
```go
//定义对外初始化日志
func InitLogger(level LogLevel,chanSize int,serviceName string)  {
	if chanSize <= 0{
		chanSize = DefaultLogChanSize;
	}

	initLogger(level,chanSize,serviceName)
}
```
至此logs.go文件修改如下：
```go
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
```
接下来开始写产生日志的相关方法，日志最重要的一个功能就是显示日志是哪个程序文件的哪行代码产生的，所以我们先写一个获取生成日志的文件名的方法。
创建util.go文件，在里面创建GetLineInfo方法来获取产生日志文件的文件名和行数。
```go
package zy_logs

import (
	"runtime"
)

//获取生成日志的文件名和行数
func GetLineInfo() (fileName string,lineNo int) {
	_, fileName, lineNo, _ = runtime.Caller(3)

	return
}

```

在日志中我们还用到追踪日志记录（这个追踪日志方面我只是简单的实现），定义追踪日志需要的文件trace_id.go
在里面定义生成traceId方法，将traceId存入上下文方法，从上下文中获取traceId方法。

```go
package zy_logs

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const (
	MaxTraceId = 100000000
)

type traceIdKey struct {}

func init()  {
	rand.Seed(time.Now().UnixNano())
}
/*获取traceId*/
func GetTraceId(ctx context.Context) (traceId string)  {
	traceId,ok := ctx.Value(traceIdKey{}).(string)
	if !ok {
		traceId = "-"
	}
	return
}
/*生成traceId*/
func GenTraceId() (traceId string) {
	now := time.Now()
	traceId = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%08d", now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), rand.Int31n(MaxTraceId))
	return
}
/*将traceId放入ctx*/
func WithTraceId(ctx context.Context, traceId string) context.Context{
	return context.WithValue(ctx, traceIdKey{}, traceId)
}
```

还可以定义一些日志分级的常量，这样就可以选择日志分级常量而不用手动输入参数了。
在util.go 中添加如下代码

```go
/*获取日志等级字符串*/
...

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

...


```
另外我们还需要定义一些方法，将日志的字段进行一些处理
```go

...

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

...

```
当日志输出到控制台的时候，希望日志的显示根据分级的不同而显示不同的颜色，起到让错误日志显眼的作用，例如颜色为红色的是严重错误的日志
首先需要在常量文件中定义一些颜色常量

```go

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

```
在util中定义
```go

...

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

...

```
新建文件color.go
```go
package zy_logs

import "fmt"

type Color uint8

func (c Color)WithColor(s string) string{
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m",uint8(c),s)
}

```
日志存储的日志文件我们是按照时间来进行切分的，定义切分日志文件的方法
同样在常量文件中定义文件切割符常量

```go

const (
	YearSeg LogFileSeg = iota
	MonthSeg
	WeekSeg
	DaySeg
	HourSeg
)

```

在util中定义
```go

...

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

...

```
访问日志有时候需要用户自己在日志数据结构中定义一些自己的字段，将用户访问日志自定义字段相关的方法写在kvs.go文件中

```go
package zy_logs

import (
	"context"
	"sync"
)

var (
	initFields sync.Once
)

type KeyVal struct {
	key interface{}
	val interface{}
}

type LogField struct {
	kvs []KeyVal
	fieldLock sync.Mutex //添加一个锁，为何加锁，因为AddField可能是多个线程在调用的
}

type kvsIdKey struct {}

/*将用户传入字段格式化*/
func (l*LogField) AddField(key , val interface{}) {
	l.fieldLock.Lock()
	l.kvs = append(l.kvs,KeyVal{key:key,val:val})
	l.fieldLock.Unlock()
}

/* 将字段存入上下文*/
func WithFieldContext(ctx context.Context) context.Context {
	fields := &LogField{}
	return context.WithValue(ctx, kvsIdKey{},fields)
}
/*向日志数据中添加其他字段*/
func AddField(ctx context.Context,key string,val interface{})  {
	field := getFields(ctx)
	if field == nil {
		return
	}
	field.AddField(key,val)
}
/*从上下文中获取其他字段*/
func getFields(ctx context.Context) *LogField {
	field, ok := ctx.Value(kvsIdKey{}).(*LogField)
	if !ok {
		return nil
	}
	return field
}

```


 
 至此这一章节的内容结束，[下个章节][第四章]主要来实现日志库的输出器的实现
 
 # 目录
 
 - [第一章 需求分析][第一章]
 - [第二章 原型实现-定义数据结构][第二章]
 - [第三章 原型实现-实现日志公共方法][第三章]
 - [第四章 原型实现-实现日志输出器][第四章]
 - [第五章 易用性封装][第五章]
 - [第六章 功能优化][第六章]
 - [第七章 功能测试][第七章]
 
 [第一章]: ../part1
 [第二章]: ../part2
 [第三章]: ../part3
 [第四章]: ../part4
 [第五章]: ../part5
 [第六章]: ../part6
 [第七章]: ../part7





