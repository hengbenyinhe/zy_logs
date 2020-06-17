package zy_logs

import (
	"runtime"
)

//获取生成日志的文件名和行数
func GetLineInfo() (fileName string,lineNo int) {
	_, fileName, lineNo, _ = runtime.Caller(3)

	return
}

