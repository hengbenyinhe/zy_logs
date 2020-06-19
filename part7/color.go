package zy_logs

import "fmt"

type Color uint8

func (c Color)WithColor(s string) string{
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m",uint8(c),s)
}
