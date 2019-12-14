# 原型实现-实现日志库函数

在[上一个章节][第二章]已经定义了一些要用到的数据结构，如果之后的开发过程中可能还需要定义数据结构，到时候会具体讲解的

这章节就实现日志库需要用到的函数或者方法，首先可以肯定的是对外暴露的日志级别函数。

## 定义对外暴露的日志级别函数

在logs.go文件中加入如下代码
```go
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



