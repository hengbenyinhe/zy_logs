package zy_logs

import (
	"context"
	"sync"
)

var (
	initFields sync.Once
)

type KeyVal struct {
	key interface{}
	val interface{}
}

type LogField struct {
	kvs []KeyVal
	fieldLock sync.Mutex //添加一个锁，为何加锁，因为AddField可能是多个线程在调用的
}

type kvsIdKey struct {}

/*将用户传入字段格式化*/
func (l*LogField) AddField(key , val interface{}) {
	l.fieldLock.Lock()
	l.kvs = append(l.kvs,KeyVal{key:key,val:val})
	l.fieldLock.Unlock()
}

/* 将字段存入上下文*/
func WithFieldContext(ctx context.Context) context.Context {
	fields := &LogField{}
	return context.WithValue(ctx, kvsIdKey{},fields)
}
/*向日志数据中添加其他字段*/
func AddField(ctx context.Context,key string,val interface{})  {
	field := getFields(ctx)
	if field == nil {
		return
	}
	field.AddField(key,val)
}
/*从上下文中获取其他字段*/
func getFields(ctx context.Context) *LogField {
	field, ok := ctx.Value(kvsIdKey{}).(*LogField)
	if !ok {
		return nil
	}
	return field
}
