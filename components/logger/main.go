package logger

import "fmt"

func Write(text string) string {
	_ = fmt.Sprintf("wow")
	return fmt.Sprintf("Log: %s", text)
}
