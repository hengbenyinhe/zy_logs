package zy_logs

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const (
	MaxTraceId = 100000000
)

type traceIdKey struct {}

func init()  {
	rand.Seed(time.Now().UnixNano())
}

/*获取traceId*/
func GetTraceId(ctx context.Context) (traceId string)  {
	traceId,ok := ctx.Value(traceIdKey{}).(string)
	if !ok {
		traceId = "-"
	}
	return
}
/*生成traceId*/
func GenTraceId() (traceId string) {
	now := time.Now()
	traceId = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%08d", now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), rand.Int31n(MaxTraceId))
	return
}
/*将traceId放入ctx*/
func WithTraceId(ctx context.Context, traceId string) context.Context{
	return context.WithValue(ctx, traceIdKey{}, traceId)
}