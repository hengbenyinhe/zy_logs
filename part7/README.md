# 功能测试

上一章节u已经基本实现日志库了，这一章主要是对功能进行测试
在logs_test中进行测试日志库暴露的方法
```go
package zy_logs

import (
	"context"
	"testing"
)

func TestFileLogger(t *testing.T) {
	//outputer, err := NewFileOutputer("logs/test.log","year")
	outputer, err := NewFileOutputer("logs/test.log",YearSeg)
	if err != nil {
		t.Errorf("init file outputer failed, err:%v", err)
		return
	}

	InitLogger(LogLevelDebug, 10000, "account")
	AddOutputer(outputer)

	Debug(context.Background(), "this is a good test")
	Trace(context.Background(), "this is a good test")
	Info(context.Background(), "this is a good test")
	Access(context.Background(), "this is a good test")
	Warn(context.Background(), "this is a good test")
	Error(context.Background(), "this is a good test")
	Stop()
}

func TestConsoleLogger(t *testing.T) {
	ctx := context.Background()
	ctx = WithFieldContext(ctx)
	ctx = WithTraceId(ctx, GenTraceId())

	AddField(ctx, "user_id", 83332232)
	AddField(ctx, "name", "kswss")



	Access(ctx, "this is a good test")

	Debug(ctx, "this is a good test")
	Trace(ctx, "this is a good test")
	Info(ctx, "this is a good test")
	Warn(ctx, "this is a good test")
	Error(ctx, "this is a good test")
	Stop()
}

```

测试结果如下：

```
hengben@ubuntu:/home/go/src/zy_logs/part7$ go test
19069-60-66 67:167:40.081 ACCESS account logs_test.go:38 2020061901164034456231 user_id=83332232 name=kswss this is a good test
19069-60-66 67:167:40.081 DEBUG account logs_test.go:40 2020061901164034456231 this is a good test
19069-60-66 67:167:40.081 TRACE account logs_test.go:41 2020061901164034456231 this is a good test
19069-60-66 67:167:40.081 INFO account logs_test.go:42 2020061901164034456231 this is a good test
19069-60-66 67:167:40.081 WARN account logs_test.go:43 2020061901164034456231 this is a good test
19069-60-66 67:167:40.081 ERROR account logs_test.go:44 2020061901164034456231 this is a good test
PASS
ok      zy_logs 0.002s
```



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