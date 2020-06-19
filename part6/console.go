package zy_logs

import "os"

type ConsoleOutputer struct {
}
//创建控制台输出实例
func NewConsoleOutputer() (log Outputer) {
	log = &ConsoleOutputer{}
	return
}
//将日志写道控制台
func (c *ConsoleOutputer)Write(data *LogData) {
	color := getLevelColor(data.level)
	text := color.WithColor(string(data.Bytes()))
	os.Stdout.Write([]byte(text))
}
//关闭写资源，控制台无任何资源需要释放所以方法里无需写任何代码
func (c *ConsoleOutputer)Close()  {

}