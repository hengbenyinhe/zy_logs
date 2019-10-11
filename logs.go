package zy_logs

import (
	"context"
)

//对外暴露日志API接口,将日志分为访问，
func Access(ctx context.Context, format string, args ...interface{}) {

}
func Debug(ctx context.Context, format string, args ...interface{}) {

}

func Trace(ctx context.Context, format string, args ...interface{}) {

}

func Info(ctx context.Context, format string, args ...interface{}) {

}

func Warn(ctx context.Context, format string, args ...interface{}) {

}

func Error(ctx context.Context, format string, args ...interface{}) {

}