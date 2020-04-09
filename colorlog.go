package colorlog

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type LogLevelEnum int
type ColorEnum int

const (
	// Log Levels
	Lnone  LogLevelEnum = 0
	Ldebug LogLevelEnum = 10
	Linfo  LogLevelEnum = 11
	Lwarn  LogLevelEnum = 12
	Lerror LogLevelEnum = 13
	Lfatal LogLevelEnum = 14

	// Colors
	Black   ColorEnum = 30
	Red     ColorEnum = 31
	Green   ColorEnum = 32
	Yellow  ColorEnum = 93
	Blue    ColorEnum = 34
	Magenta ColorEnum = 35
	Cyan    ColorEnum = 36
	Grey    ColorEnum = 37
	White   ColorEnum = 97

	// Formats
	screenFormat string = "\033[0;%dm[%s] <%v> %s\033[0m"
	colorFormat  string = "\033[0;%dm%s\033[0m"
	fileFormat   string = "[%s] <%v> %s"
)

func (ll LogLevelEnum) String() string {
	switch ll {
	case Lnone:
		return " NONE"
	case Ldebug:
		return "DEBUG"
	case Linfo:
		return " INFO"
	case Lwarn:
		return " WARN"
	case Lerror:
		return "ERROR"
	case Lfatal:
		return "FATAL"
	default:
		return fmt.Sprintf("%d", ll)
	}
}

type ColorLog struct {
	screenFormat string
	fileFormat   string
	LogLevel     LogLevelEnum
	logFile      *os.File
	logWriter    *bufio.Writer
	isFileLogger bool
}

func New(level LogLevelEnum) *ColorLog {
	// default to INFO if level requested is beyond the valid enums
	if level > 14 {
		level = 11
	}
	return &ColorLog{LogLevel: LogLevelEnum(level), screenFormat: screenFormat, fileFormat: fileFormat, isFileLogger: false}
}

func NewFileLog(level LogLevelEnum, filename string) *ColorLog {
	if level > 14 {
		level = 11
	}

	logfile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(logfile)
	return &ColorLog{LogLevel: LogLevelEnum(level),
		screenFormat: screenFormat,
		fileFormat:   fileFormat,
		logFile:      logfile,
		logWriter:    writer,
		isFileLogger: true}
}

func (l *ColorLog) Debug(msg string) {
	l.Print(msg, 10, Green)

	if l.isFileLogger {
		l.Write(msg, 10)
	}
}

func (l *ColorLog) Info(msg string) {
	l.Print(msg, 11, Grey)

	if l.isFileLogger {
		l.Write(msg, 11)
	}
}

func (l *ColorLog) Warn(msg string) {
	l.Print(msg, 12, Yellow)

	if l.isFileLogger {
		l.Write(msg, 12)
	}
}

func (l *ColorLog) Error(msg string) {
	l.Print(msg, 13, Red)

	if l.isFileLogger {
		l.Write(msg, 13)
	}
}

func (l *ColorLog) Fatal(msg string) {
	l.Print(msg, 14, Magenta)

	if l.isFileLogger {
		l.Write(msg, 14)
	}
}

func (l *ColorLog) Print(msg string, level LogLevelEnum, color ColorEnum) {
	if l.LogLevel > 0 && LogLevelEnum(level) >= l.LogLevel {
		fmt.Printf(l.screenFormat, color, l.TimeStamp(), LogLevelEnum(level), msg)
	}
}

func (l *ColorLog) Printc(msg string, level LogLevelEnum, color ColorEnum) {
	if l.LogLevel > 0 && level >= l.LogLevel {
		fmt.Printf(colorFormat, color, msg)
	}
}

func (l *ColorLog) Write(msg string, level LogLevelEnum) {
	if l.isFileLogger && l.LogLevel > 0 && LogLevelEnum(level) >= l.LogLevel {
		str := fmt.Sprintf(l.fileFormat, l.TimeStamp(), LogLevelEnum(level), msg)
		l.logWriter.WriteString(str)
		l.logWriter.Flush()
	}
}

func (l *ColorLog) Close() {
	if l.isFileLogger {
		l.logWriter.Flush()
		l.logFile.Close()
	}
}

func (l *ColorLog) TimeStamp() string {
	return time.Now().Format(time.RFC3339)
}

func (l *ColorLog) SetLogLevel(level LogLevelEnum) {
	l.LogLevel = level
}
