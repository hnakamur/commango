package executil

import (
	"fmt"
	"strings"
	"time"
)

type LogLevel int

const (
	Info LogLevel = iota
	Err
)

func (level LogLevel) String() string {
	switch level {
	case Info:
		return "info"
	case Err:
		return "err"
	default:
		panic("not reached")
	}
}

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

func Logf(level LogLevel, format string, args... interface{}) {
	header := fmt.Sprintf("time:%s\tlevel:%s",
		time.Now().Format(RFC3339Milli),
		level)
	lines := strings.Split(
		strings.TrimRight(fmt.Sprintf(format, args...), "\n"),
		"\n")
	for _, line := range(lines) {
		if level == Info {
			fmt.Printf("%s\tmessage:%s\n", header, line)
		} else if level == Err {
			fmt.Printf("\x1b[31;1m%s\tmessage:%s\x1b[0m\n", header, line)
		}
	}
}

type Logger struct {
	level LogLevel
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{level}
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.Logf("%s", string(p))
	return len(p), nil
}

func (l *Logger) Logf(format string, args... interface{}) {
	Logf(l.level, format, args...)
}
