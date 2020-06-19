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
