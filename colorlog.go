package colorlog

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// LogLevelEnum : logging levels
type LogLevelEnum int

// ColorEnum : color keys
type ColorEnum int

const (
	// Log Levels

	Lnone  LogLevelEnum = 0  // Log nothing
	Ldebug LogLevelEnum = 10 // Log Debug and above
	Linfo  LogLevelEnum = 11 // Log Info and above
	Lwarn  LogLevelEnum = 12 // Log Warn and above
	Lerror LogLevelEnum = 13 // Log Error and above
	Lfatal LogLevelEnum = 14 // Log Fatal

	// Colors

	Black   ColorEnum = 30 // Black
	Red     ColorEnum = 31 // Red
	Green   ColorEnum = 32 // Green
	Yellow  ColorEnum = 93 // Yellow
	Blue    ColorEnum = 34 // Blue
	Magenta ColorEnum = 35 // Magenta
	Cyan    ColorEnum = 36 // Cyan
	Grey    ColorEnum = 37 // Grey
	White   ColorEnum = 97 // White

	// Formats

	screenFormat string = "\033[0;%dm[%s] <%v> %s\033[0m" // Format to use for colorized log lines
	colorFormat  string = "\033[0;%dm%s\033[0m"           // Format to use for colorizing appends to log lines
	fileFormat   string = "[%s] <%v> %s"                  // Format to use for log file lines (same as screenFormat minus the color
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

// ColorLog structure
type ColorLog struct {
	screenFormat string
	fileFormat   string
	LogLevel     LogLevelEnum
	logFile      *os.File
	logWriter    *bufio.Writer
	isFileLogger bool
	noColor      bool
}

// New : Create a new ColorLog
func New(level LogLevelEnum) *ColorLog {
	if level > 14 {
		level = 11
	}
	return &ColorLog{LogLevel: LogLevelEnum(level), screenFormat: screenFormat, fileFormat: fileFormat, isFileLogger: false, noColor: false}
}

// NewColorless : Create a new colorless ColorLog
func NewColorless(level LogLevelEnum) *ColorLog {
	if level > 14 {
		level = 11
	}
	return &ColorLog{LogLevel: LogLevelEnum(level), screenFormat: screenFormat, fileFormat: fileFormat, isFileLogger: false, noColor: true}
}

// NewFileLog : Create a new ColorLog file logger
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

// Debug log message
func (l *ColorLog) Debug(msg string, opts ...interface{}) {
	msg = fmt.Sprintf(msg, opts...)

	if !l.noColor {
		l.Print(msg, 10, Green)
	} else {
		fmt.Printf(l.fileFormat, l.TimeStamp(), LogLevelEnum(10), msg)
	}

	if l.isFileLogger {
		l.Write(msg, 10)
	}
}

// Info log message
func (l *ColorLog) Info(msg string, opts ...interface{}) {
	msg = fmt.Sprintf(msg, opts...)

	if !l.noColor {
		l.Print(msg, 11, Grey)
	} else {
		fmt.Printf(l.fileFormat, l.TimeStamp(), LogLevelEnum(11), msg)
	}

	if l.isFileLogger {
		l.Write(msg, 11)
	}
}

// Warn log message
func (l *ColorLog) Warn(msg string, opts ...interface{}) {
	msg = fmt.Sprintf(msg, opts...)

	if !l.noColor {
		l.Print(msg, 12, Yellow)
	} else {
		fmt.Printf(l.fileFormat, l.TimeStamp(), LogLevelEnum(12), msg)
	}

	if l.isFileLogger {
		l.Write(msg, 12)
	}
}

// Error log message
func (l *ColorLog) Error(msg string, opts ...interface{}) {
	msg = fmt.Sprintf(msg, opts...)

	if !l.noColor {
		l.Print(msg, 13, Red)
	} else {
		fmt.Printf(l.fileFormat, l.TimeStamp(), LogLevelEnum(13), msg)
	}

	if l.isFileLogger {
		l.Write(msg, 13)
	}
}

// Fatal log message
func (l *ColorLog) Fatal(msg string, opts ...interface{}) {
	msg = fmt.Sprintf(msg, opts...)

	if !l.noColor {
		l.Print(msg, 14, Magenta)
	} else {
		fmt.Printf(l.fileFormat, l.TimeStamp(), LogLevelEnum(14), msg)
	}

	if l.isFileLogger {
		l.Write(msg, 14)
	}
}

// Print : Print helper function
func (l *ColorLog) Print(msg string, level LogLevelEnum, color ColorEnum) {
	if l.LogLevel > 0 && LogLevelEnum(level) >= l.LogLevel {
		fmt.Printf(l.screenFormat, color, l.TimeStamp(), LogLevelEnum(level), msg)
	}
}

// Printc : Print colorized message
func (l *ColorLog) Printc(level LogLevelEnum, color ColorEnum, msg string, opts ...interface{}) {
	msg = fmt.Sprintf(msg, opts...)
	if l.LogLevel > 0 && level >= l.LogLevel {
		fmt.Printf(colorFormat, color, msg)
	}
}

// Write : Logfile write helper function
func (l *ColorLog) Write(msg string, level LogLevelEnum) {
	if l.isFileLogger && l.LogLevel > 0 && LogLevelEnum(level) >= l.LogLevel {
		str := fmt.Sprintf(l.fileFormat, l.TimeStamp(), LogLevelEnum(level), msg)
		l.logWriter.WriteString(str)
		l.logWriter.Flush()
	}
}

// Close : Helper function to close a file logger
func (l *ColorLog) Close() {
	if l.isFileLogger {
		l.logWriter.Flush()
		l.logFile.Close()
	}
}

// TimeStamp : Get the current Timestamp
func (l *ColorLog) TimeStamp() string {
	return time.Now().Format(time.RFC3339)
}

// SetLogLevel : Set the logger's log level
func (l *ColorLog) SetLogLevel(level LogLevelEnum) {
	l.LogLevel = level
}

// GetLogLevel : Return the logger's log level
func (l *ColorLog) GetLogLevel() LogLevelEnum {
	return l.LogLevel
}
