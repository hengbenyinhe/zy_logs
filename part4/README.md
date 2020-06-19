# 原型实现-实现日志输出器

在[上一个章节][第三章]我们定义了一些日志库需要用到的方法，这一章节主要来实现日志输出器。
因为这个日志库主要输出日志到日志文件和控制台，默认日志输出器为控制台输出器。我们先来实现文件输出器。
定义一个日志文件输出器。新建一个file.go文件。

定义一个输出器结构体Outputer

```go

import (
	"fmt"
	"os"
	"path/filepath"
)

type Outputer interface {
	
}

```
因为日志文件是按照时间进行切分的，所以需要一个文件输出器操作状态记录机构

```go
type FileOutputerOptions struct {
	filename string
	lastSplitHour int
}
```
接下来定义文件输出器机构体
```go

type FileOutputer struct {
	file           *os.File
	accessFile     *os.File
	option         *FileOutputerOptions
}

```

定义初始化文件输出器方法
```go

func (f *FileOutputer)initFile(filename,originFilename string) (file *os.File,err error) {
	file, err = os.OpenFile(filename,os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755
	if err != nil {
		err = fmt.Errorf("open faile %s failed, err:%v", filename, err)
		return
	}

	os.Symlink(filename,originFilename)
	return
}

```
定义生成日志文件的函数,因为访问类型的日志需要另外分日志文件，所以也需要定义一个生成访问日志文件的方法
```go
func (f *FileOutputer)getCurFilename() (curFilename,originFilename string)  {
	now := time.Now()
	curFilename = fmt.Sprintf("%s.%04d%02d%02d%02d",f.option.filename,now.Year(),now.Month(),
		now.Day(),now.Hour())
	originFilename = f.option.filename
	return
}

func (f *FileOutputer)getCurAccessFilename() (curAccessFilename,originAccessFilename string) {
	now := time.Now()
	curAccessFilename = fmt.Sprintf("%s.access.%04d%02d%02d%02d",f.option.filename,now.Year(),now.Month(),
		now.Day(),now.Hour())
	originAccessFilename = f.option.filename
	return
	
}
```
初始化函数
```go

func (f *FileOutputer)init() (err error) {
	curFilename, originFileame := f.getCurFilename()
	f.file,err = f.initFile(curFilename,originFileame)
	if err !=nil {
		return
	}

	accessFilename, originAccessFilename := f.getCurAccessFilename()
	f.accessFile, err = f.initFile(accessFilename, originAccessFilename)
	if err !=nil {
		return
	}

	f.option.lastSplitHour = time.Now().Hour()
	return
}

```
在输出日志到文件的时候，我们需要判断是不是要切换日志文件，定义一个检查日志文件切换的方法
```go

func (f *FileOutputer)checkSplitFile(curTime time.Time) {
	hour := curTime.Hour()
	if hour == f.option.lastSplitHour{
		return
	}

	f.init()
}

```
日志输出器关键是要有写方法和关闭输出器的方法
```go
func (f *FileOutputer)Write(data *LogData) {
	f.checkSplitFile(data.curTime)
	file := f.file
	if data.level == LogLevelAccess{
		file = f.accessFile
	}

	file.Write(data.Bytes())
}

func (f *FileOutputer)Close() {

	f.file.Close()
	f.accessFile.Close()
	
}
```
以上就基本完成文件输出器的内容了，那我们再来定义控制台输出器模块，控制台输出器默认为日志库的日志输出器。
新建console.go文件,i前面提到过输出器中主要包含Write和Close这两个方法，还有实例化控制台输出器的方法NewConsoleOutputer
先定义写方法
```go
func (c *ConsoleOutputer)Write(data *LogData) {
	color := getLevelColor(data.level)
	text := color.WithColor(string(data.Bytes()))
	os.Stdout.Write([]byte(text))
}
```
关闭方法因为控制台无任何资源需要释放所以方法里面无需写任何代码
```go
func (c *ConsoleOutputer)Close()  {

}
```
定义实例化方法
```go
func NewConsoleOutputer() (log Outputer) {
	log = &ConsoleOutputer{}
	return
}

```
[下个章节][第五章]主要来实现日志库的记录器以及对日志库进行封装

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