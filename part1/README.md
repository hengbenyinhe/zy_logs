# 需求分析

自己是在用go微服务框架go-micro做项目时。需要将操作数据的操作写入到日志文件中，然micro提供的日志包只支持输出到控制台且不支持日志分等级。
当时自己就想着网上查找一下资料，怎么自己简单写一个满足自己需求的日志库。<br />

一般的情况下，日志库支持文件写入和console显示，有的甚至支持写入到网络服务中，但是我就简单实现写入到日志文件和输出到控制台的时候根据日志等级显示不同颜色就可以了。

## 定义日志的打印级别

zy_log日志包提供6种级别的日志，当日志输出到控制台时根据日志等级来显示相应的颜色。

1. Debug: 调试程序，日志最详细。但是会影响程序的性能。
2. Trace: 追踪问题。
3. Info: 打印日志运行中比较重要的信息，比如状态变化日志。
4. Warn: 警告，说明程序中出现了潜在的问题。
5. Error: 错误，程序运行发生了错误。
6. Access：访问，用户对程序的接口访问日志。

## 日志存储位置

支持两种日志输出方式

1. 直接输出到console
2. Trace: 追踪问题。

## 日志数据结构

日志产生日期+日志产生时间（秒）+日志级别+产生日志服务名+产生日志的文件名+行号+traceId+[其他字段…]+日志信息

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