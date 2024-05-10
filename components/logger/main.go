package logger

import ("fmt")

func Write(text string) string {
	return fmt.Sprintf("Log: %s", text)
}

