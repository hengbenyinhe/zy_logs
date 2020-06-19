# 

这个章节主要是对文件输出器进一步进行完善，日志文件的切分有时候根据天划分，有时根据月划分，所以我们提供可以给用户自己定义时间切分日志文件的功能。
而且要考虑到当切换日志的状态记录，比如万一第一天日志产生时间为11点，第二天产生的时间也为11点，那么这个时候根据原来逻辑则不会切换文件。
所以将file文件做如下修改：
```go
package zy_logs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileOutputerOptions struct {
	filename   string //文件名
	lastSplitHour int //记录上次切换文件小时，初始设计时用到。现在已经不用
	lastSplitTime string //记录上次切换文件时间，考虑这种情况比如第一天日志产生时间为11点，第二天产生的时间也为11点，那么这个时候根据原来逻辑则不会切换文件
	lastSplitWeek time.Time //便于以周为切分周期的情况
}

type FileOutputer struct {
	file    *os.File
	accessFile *os.File
	option   *FileOutputerOptions
}

var (
	cutTime string //文件切分时段
)

//创建文件输出器实例
func NewFileOutputer(filename string,cuttime LogFileSeg) (Outputer,error) {
	cutTime = getSegText(cuttime)

	filename,err := filepath.Abs(filename)
	if err != nil {
		return nil,err
	}
	option := &FileOutputerOptions{
		filename:filename,
	}

	log := &FileOutputer{
		option:option,
	}
	err = log.init()
	return log,err
}
/*生成日志文件名*/
func (f *FileOutputer)getCurFilename(cuttime string) (curFilename,originFilename string) {
	now := time.Now()
	/*curFilename = fmt.Sprintf("%s.%04d%02d%02d%02d",f.option.filename,now.Year(),now.Month(),
	now.Day(),now.Hour())*/
	switch cuttime {
	case "day":
		curFilename = fmt.Sprintf("%s.%04d%02d%02d",f.option.filename,now.Year(),now.Month(),now.Day())
	case "week":
		curFilename = fmt.Sprintf("%s.%04d%02d%02d",f.option.filename,now.Year(),now.Month(),now.Day())
	case "month":
		curFilename = fmt.Sprintf("%s.%04d%02d",f.option.filename,now.Year(),now.Month())
	case "year":
		curFilename = fmt.Sprintf("%s.%04d",f.option.filename,now.Year())
	default:
		curFilename = fmt.Sprintf("%s.%04d%02d%02d%02d",f.option.filename,now.Year(),now.Month(),
			now.Day(),now.Hour())
	}
	originFilename = f.option.filename
	return
}
/*生成访问日志文件名*/
func (f *FileOutputer)getCurAccessFilename(cuttime string) (curAccessFilename,originAccessFilename string) {
	now := time.Now()
	/*curAccessFilename = fmt.Sprintf("%s.access.%04d%02d%02d%02d",f.option.filename,now.Year(),now.Month(),
	now.Day(),now.Hour())*/
	switch cuttime {
	case "day":
		curAccessFilename = fmt.Sprintf("%s.access.%04d%02d%02d",f.option.filename,now.Year(),now.Month(),now.Day())
	case "week":
		curAccessFilename = fmt.Sprintf("%s.access.%04d%02d%02d",f.option.filename,now.Year(),now.Month(),now.Day())
	case "month":
		curAccessFilename = fmt.Sprintf("%s.access.%04d%02d",f.option.filename,now.Year(),now.Month())
	case "year":
		curAccessFilename = fmt.Sprintf("%s.access.%04d",f.option.filename,now.Year())
	default:
		curAccessFilename = fmt.Sprintf("%s.access.%04d%02d%02d%02d",f.option.filename,now.Year(),now.Month(),
			now.Day(),now.Hour())
	}
	originAccessFilename = fmt.Sprintf("%s.acccess", f.option.filename)
	return
}

/*初始化文件*/
func (f *FileOutputer)initFile(filename,originFilename string) (file *os.File,err error) {
	file, err = os.OpenFile(filename,os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		err = fmt.Errorf("open faile %s failed, err:%v", filename, err)
		return
	}
	os.Symlink(filename,originFilename)
	return
}
/*初始化日志输出器文件*/
func (f *FileOutputer)init() (err error) {
	//curFilename, originFileame := f.getCurFilename()
	curFilename, originFileame := f.getCurFilename(cutTime)
	f.file,err = f.initFile(curFilename,originFileame)
	if err !=nil {
		return
	}

	accessFilename, originAccessFilename := f.getCurAccessFilename(cutTime)
	f.accessFile, err = f.initFile(accessFilename, originAccessFilename)
	if err != nil {
		return
	}

	//f.option.lastSplitHour = time.Now().Hour()
	switch cutTime {
	case "day":
		f.option.lastSplitTime = fmt.Sprintf("%04d%02d%02d",time.Now().Year(),time.Now().Month(), time.Now().Day())
	case "week":
		f.option.lastSplitWeek = time.Now()
	case "month":
		f.option.lastSplitTime = fmt.Sprintf("%04d%02d",time.Now().Year(),time.Now().Month())
	case "year":
		f.option.lastSplitTime = fmt.Sprintf("%04d",time.Now().Year())
	default:
		f.option.lastSplitTime = fmt.Sprintf("%04d%02d%02d%02d",time.Now().Year(),time.Now().Month(), time.Now().Day(),time.Now().Hour())
	}
	return

}
/*判断是否切换文件*/
func (f *FileOutputer)checkSplitFile(curTime time.Time) {
	var date string

	switch cutTime {
	case "day":
		date = fmt.Sprintf("%04d%02d%02d",curTime.Year(),curTime.Month(), curTime.Day())
	case "month":
		date = fmt.Sprintf("%04d%02d",curTime.Year(),curTime.Month())
	case "year":
		date = fmt.Sprintf("%04d",curTime.Year())
	default:
		date = fmt.Sprintf("%04d%02d%02d%02d",time.Now().Year(),time.Now().Month(), time.Now().Day(),time.Now().Hour())
	}

	if cutTime =="week" {
		subM := curTime.Sub(f.option.lastSplitWeek)
		timeHour:=int(subM.Hours())
		days := timeHour/24
		if days<7 {
			return
		}
	}else {
		if date == f.option.lastSplitTime{
			return
		}
	}

	f.init()
}
/*写入日志到文件*/
func (f *FileOutputer)Write(data *LogData)  {
	f.checkSplitFile(data.curTime)
	file := f.file
	if data.level == LogLevelAccess{
		file = f.accessFile
	}

	file.Write(data.Bytes())
}
/*关闭，释放资源*/
func (f *FileOutputer)Close() {
	f.file.Close()
	f.accessFile.Close()
}


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