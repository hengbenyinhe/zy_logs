# 原型实现

根据上一个章节的[需求分析][第一章]知道zy_log日志库分为六个日志等级

## 初始化日志需要用到的数据结构

使用日志库时是不一定需要初始化的，如果没有初始化，则服务名等数据则为默认值。因为日志等级提供了六种固定的等级，所以在初始化的时候可以
InitLogger
接下来开始开发，创建项目入口文件logs.go文件，项目采用go module进行包管理，在终端执行go mod init zy_logs命令，生成go.mod文件。在logs.go里面定义一个日志等级数据类型LogLevel，以便之后定义日志等级属性;同时此时也可以定义日志库需要用到的数据
```go
    package zy_logs
    
    import (
    	"sync"
    	"time"
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
```

所以需要定义日志等级属性如下：
这么做的目的主要是用户不用用户去输入日志等级常量。在项目中创建consts.go文件

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
```
至此这一章节的内容结束，[下个章节][第三章]将定义日志需要用到的接口

# 目录

- [第一章 需求分析][第一章]
- [第二章 原型实现-定义数据结构][第二章]
- [第二章 原型实现-实现日志接口][第三章]
- [第三章 易用性封装][第四章]
- [第四章 功能优化][第五章]
- [第五章 功能测试][第六章]
- [第六章 学习总结][第七章]

[第一章]: ../part1
[第二章]: ../part2
[第三章]: ../part3
[第四章]: ../part4
[第五章]: ../part5
[第六章]: ../part6
[第七章]: ../part7