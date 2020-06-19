
#zy_log日志包组件
>golang 编写的简单日志包

##日志等级
zy_log日志包提供6种级别的日志，当日志输出到控制台时根据日志等级来显示相应的颜色。<br/>
&emsp;&emsp;1. Debug: 调试程序，日志最详细。但是会影响程序的性能。<br/>
&emsp;&emsp;2. Trace: 追踪问题。<br/>
&emsp;&emsp;3. Info: 打印日志运行中比较重要的信息，比如状态变化日志。<br/>
&emsp;&emsp;4. Warn: 警告，说明程序中出现了潜在的问题。<br/>
&emsp;&emsp;5. Error: 错误，程序运行发生了错误。<br/>
&emsp;&emsp;6. Access：访问，用户对程序的接口访问日志。<br/>
##日志输出方式
支持两种日志输出方式，输出到控制台时不需要初始化输出器，但是输出到文件就必须初始化日志输出器。输出到文件支持以年或月或周或日或小时进行日志切分的,且是相对与程序main文件的同级<br/>
&emsp;&emsp;1. 输出到控制台<br/>
&emsp;&emsp;2. 输出到文件<br/>
##日志数据结构
日志产生日期+日志产生时间（秒）+日志级别+产生日志服务名+产生日志的文件名+行号+traceId+[其他字段...]+日志信息<br/>
##组件API
|名称|说明|参数|
|:----:|----|----|
| InitLogger  | 初始化日志，主要是初始化日志的一些基本数据,如果没有初始化则日志数据将会为默认值 |1.level:日志等级，这个参数需要用到LogLevel属性的值;2.chanSize:日志输出管道大小;3.serviceName:产生日志服务名 |
| Access  | 访问日志记录器，生成访问日志数据并且将日志写入到输出器 | 1.ctx:上下文;2.format:日志信息输出格式;3.args：日志数据参数|
| Debug  | 调试日志记录器，生成访问日志数据并且将日志写入到输出器 | 1.ctx:上下文;2.format:日志信息输出格式;3.args：日志数据参数 |
| Trace  | 追踪日志记录器，生成访问日志数据并且将日志写入到输出器 | 1.ctx:上下文;2.format:日志信息输出格式;3.args：日志数据参数 |
| Info  | 信息日志记录器，生成访问日志数据并且将日志写入到输出器 | 1.ctx:上下文;2.format:日志信息输出格式;3.args：日志数据参数 |
| Warn  | 警告日志记录器，生成访问日志数据并且将日志写入到输出器 | 1.ctx:上下文;2.format:日志信息输出格式;3.args：日志数据参数 |
| Error  | 错误日志记录器，生成访问日志数据并且将日志写入到输出器 | 1.ctx:上下文;2.format:日志信息输出格式;3.args：日志数据参数 |
| Stop  | 停止日志输出器，释放输出日志资源，也可以用于输出器的切换 | 无 |
| NewFileOutputer  | 生成文件输出器实例 | 1.filename:文件名（包含文件所在相对路径）2.cuttime: 文件切分时段，这个参数要用到LogFileSeg属性的值 |
| NewConsoleOutputer  | 生成控制台输出器 | 无 |
| AddOutputer  | 添加输出器 | 1.outputer：输出器实例 |
| WithFieldContext  | 将字段写入上下文，便于之后AddField方法添加字段 | 1.ctx:上下文 |
| AddField  | 添加日志数据字段，用于访问日志，例如添加用户名之类的字段 | 1.ctx:上下文;2.key:字段名;3.val：字段值 |
| GenTraceId  | 生成traceId字符串 | 无 |
| WithTraceId  | 将traceId字符串放入上下文 | 1.ctx:上下文;2.traceId:traceId字符串 |
##组件属性
|名称|值|说明|
|:----:|----|----|
|LogLevel<br/>（日志等级属性）| LogLevelDebug,<br/>LogLevelTrace,<br/>LogLevelAccess,<br/>LogLevelInfo,<br/>LogLevelWarn,<br/>LogLevelError  | 值为LogLevelDebug时初始化debug日志等级;<br/>值为LogLevelTrace时初始化trace日志等级;<br/>值为LogLevelAccess时初始化access日志等级;<br/>值为LogLevelInfo时初始化info日志等级;<br/>值为LogLevelWarn时初始化warn日志等级;<br/>值为LogLevelError时初始化error日志等级;  |
|LogFileSeg<br/> （日志文件切分属性）| YearSeg,<br/>MonthSeg,<br/>WeekSeg,<br/>DaySeg,<br/>HourSeg  | 值为YearSeg时日志文件以年为切分时段进行切分;<br/>值为MonthSeg时日志文件以月为切分时段进行切分<br/>值为WeekSeg时日志文件以周为切分时段进行切分<br/>值为DaySeg时日志文件以天为切分时段进行切分<br/>值为HourSeg时日志文件以小时为切分时段进行切分  |
##开发环境
日志库开发所使用的golang版本为1.12.5，因此使用的时候为避免因为golang版本而导致日志包无法使用，需要安装的golang版本≥1.12.5
##使用方式
1.引入<br/>
 由于无法go get命令无法下载公司代码自托管平台的项目，所以不能采用类似github平台上一些依赖包的引入方式。
 只能先从私有仓库上将依赖包下载下来作为本地依赖包来使用。<br/>
 将下载下来的依赖包放入项目中，通过go module进行管理依赖包。<br/>
  [ 仓库地址:http://oa.zyqwt.com/lifangxiu/zy_logs](http://oa.zyqwt.com/lifangxiu/zy_logs) <br/>
2.使用例子<br/>
备注：引入的路径根据自己的实际情况而定
+ 将日志只输出到控制台<br/>
将日志输出到控制台的话，可以不需要初始化日志输出器，日志默认输出器为控制台输出器<br/>
代码如下：<br/>

```package main
  
  import (
  	"context"
  	log "workspace/config_server/plugins/zy_logs"
  )
  
  func main() {
	ctx := context.Background()
	ctx = log.WithTraceId(ctx, log.GenTraceId())     //在上下文中存入traceId
	
	ctx = log.WithFieldContext(ctx)         //将字段添加到上下文中
	log.AddField(ctx, "user_id", 83332232)      //访问日志添加用户id字段
	log.AddField(ctx, "name", "kswss")          //访问日志添加用户名字段

	log.InitLogger(log.LogLevelDebug, 10000, "nihaao")   //初始化日志组件
	//生成并记录各个等级日志
	log.Access(ctx, "日志信息格式字符串:v%", "我是参数")
	log.Debug(ctx, "日志信息格式字符串:v%","我是参数")
	log.Trace(ctx, "日志信息格式字符串:v%","我是参数")
	log.Info(ctx, "日志信息格式字符串:v%","我是参数")
	log.Warn(ctx, "日志信息格式字符串:v%","我是参数")
	log.Error(ctx, "日志信息格式字符串:v%","我是参数")
	log.Stop()    //关闭控制台输出器资源
  }
```
+ 将日志只输出到文件<br/>
注意！将日志输出到文件必须初始化输出器<br/>
代码如下：<br/>

```
  package main
  
  import (
  	"context"
  	log "workspace/config_server/plugins/zy_logs"
  )
  
  func main() {
        //初始化文件输出器
	outputer, err := log.NewFileOutputer("logs/test.log",log.YearSeg)   //在此之前请确保main.go同级存在logs目录，如果不存在则写日志文件失败，此例子是以年为切分时段
	if err != nil {
		t.Errorf("init file outputer failed, err:%v", err)
		return
	}
	log.AddOutputer(outputer)//添加文件日志输出器
	
	log.InitLogger(log.LogLevelDebug, 10000, "account")   //初始化日志组件
	//生成并记录各个等级日志
	log.Debug(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Trace(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Info(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Access(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Warn(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Error(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Stop()   //关闭控制台输出器资源
  }
```
+ 将日志同时输出到文件和控制台<br/>
注意！将日志输出到文件必须初始化输出器<br/>
代码如下：<br/>

```
  package main
  
  import (
  	"context"
  	log "workspace/config_server/plugins/zy_logs"
  )
  
  func main() {
	outputer, err := log.NewFileOutputer("logs/test.log",log.DaySeg) //此例子是以天为切分时段
	if err != nil {
		t.Errorf("init file outputer failed, err:%v", err)
		return
	}
	outputer1:= log.NewConsoleOutputer()

	log.InitLogger(log.LogLevelDebug, 10000, "account")
	log.AddOutputer(outputer)
	log.AddOutputer(outputer1)
	
	log.Debug(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Trace(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Info(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Access(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Warn(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Error(context.Background(), "日志信息格式字符串:v%", "我是参数")
	log.Stop()
  }
```


+ 将日志有的输出到文件有的输出到控制台<br/>
注意！将日志输出到文件必须初始化输出器<br/>
代码如下：<br/>

```
  package main
  
  import (
  	"context"
  	log "workspace/config_server/plugins/zy_logs"
  )
  
  func main() {
        log.Debug(ctx,"我将输出到控制台,%v",1)
        log.Stop()  //停止当前输出器，可切换输出器
        
        outputer, err := log.NewFileOutputer("logs/test.log",log.DaySeg)//此例子是以天为切分时段
        if err != nil {
            log.Error(ctx,"init file outputer failed, err:%v", err)
            return
        }
        log.InitLogger(log.LogLevelDebug, 10000, "account")
        log.AddOutputer(outputer)
        log.Info(context.Background(), "我将输出到文件,%v"，1)
        log.Stop()  //停止当前输出器，可切换输出器
        
        log.Debug(context.Background(), "我也将输出到控制台")
        log.Stop()
  }
```


  
