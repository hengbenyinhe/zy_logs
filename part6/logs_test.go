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
