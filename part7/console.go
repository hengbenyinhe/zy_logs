package zy_logs

import "os"

type ConsoleOutputer struct {

}

func NewConsoleOutputer() (log Outputer) {
	log = &ConsoleOutputer{}
	return
}

func (c *ConsoleOutputer)Write(data *LogData) {
	color := getLevelColor(data.level)
	text := color.WithColor(string(data.Bytes()))
	os.Stdout.Write([]byte(text))
}

func (c *ConsoleOutputer)Close() {

}