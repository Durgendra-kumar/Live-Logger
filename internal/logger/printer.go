package logger

import (
	"fmt"
	"time"

	internal "github.com/Durgendra-kumar/Live-Logger/internal"
)

// ANSI color codes — makes the terminal output readable at a glance
const (
	colorReset  = "\033[0m"
	colorGray   = "\033[90m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
)

// Print writes a colorized, formatted log line to stdout.
func Print(event internal.LogEvent) {
	color := colorForLevel(event.Level)
	ts := event.Timestamp.Format(time.TimeOnly) // "15:04:05"

	fmt.Printf("%s%s%s  %s%-7s%s  %s%-16s%s  %s\n",
		colorGray, ts, colorReset,
		color, event.Level, colorReset,
		colorBold, event.Service, colorReset,
		event.Message,
	)
}

func colorForLevel(l internal.Level) string {
	switch l {
	case internal.DEBUG:
		return colorGray
	case internal.INFO:
		return colorGreen
	case internal.WARN:
		return colorYellow
	case internal.ERROR:
		return colorRed
	default:
		return colorReset
	}
}
