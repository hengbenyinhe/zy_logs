# 易用性封装

封装即是对外隐藏对象的属性和实现的细节，仅对外公开接口，提高了代码的重用性,易用性和安全性。
[上一个章节][第四章]实现了日志输出器，在这一章节我们来实现日志记录器。将日志数据放入管道。
在logs.go中实现日志记录器

```go
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

```

 对外暴露添加日志输出器的方法
 ```go
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
```
 
 将日志数据放入通道之后，需要从通道里面取出日志数据
 定义从通道获取日志数据方法
 
 ```go
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
```
接下来，我们在初始化日志方法中添加一个异步的协程来完成日志的记录和输出到相应的日志输出器。
修改initLogger方法

```go
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
```
日志输出到文件的时候需要释放资源，所以在这里定义一个关闭日志资源函数,对外暴露Stop()关闭资源方法

```go

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

```

最后将记录器写入对外暴露的日志分级API接口

```go

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

```

[下个章节][第六章]主要是对日志库功能完善一下

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