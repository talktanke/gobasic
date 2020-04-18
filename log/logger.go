package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

const (
	colorRed = uint8(iota + 91)
	colorGreen
	colorYellow
	colorBlue
	colorMagenta // 洋红
	colorCyan
)

var (
	file *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC"}

	levelColors = []func(test string) string{
		func(s string) string {
			return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorBlue, s)
		},
		func(s string) string {
			return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorGreen, s)
		},
		func(s string) string {
			return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorRed, s)
		},
		func(s string) string {
			return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorYellow, s)
		},
		func(s string) string {
			return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorMagenta, s)
		},
		func(s string) string {
			return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorCyan, s)
		},
	}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
	PANIC
)

// Setup initialize the log instance
func init() {
	logger = log.New(os.Stdout, DefaultPrefix, log.LstdFlags)
}

func SetWriter(w io.Writer) {
	logger.SetOutput(w)
}

func SetOutputFilename(filename string) error {
	// 不知道这样操作能不能行？先输出到标准输出
	logger.SetOutput(os.Stdout)
	if file != nil {
		file.Close()
		file = nil
	}
	// 创建或者打开文件，然后更改叔叔
	var err error
	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	logger.SetOutput(file)
	return nil
}

// SetFlags 修改flag
func SetFlags(flag int) {
	logger.SetFlags(flag)
}

// Debug output logs at debug level
func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v...)
}

// Info output logs at info level
func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v...)
}

// Warn output logs at warn level
func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v...)
}

// Error output logs at error level
func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v...)
}

// Fatal output logs at fatal level
func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v...)
}

// Panic output logs at fatal level
func Panic(v ...interface{}) {
	setPrefix(PANIC)
	logger.Panic(v...)
}

// Debugf output logs at debug level
func Debugf(s string, v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(fmt.Sprintf(s, v...))
}

// Infof output logs at info level
func Infof(s string, v ...interface{}) {
	setPrefix(INFO)
	logger.Println(fmt.Sprintf(s, v...))
}

// Warnf output logs at warn level
func Warnf(s string, v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(fmt.Sprintf(s, v...))
}

// Errorf output logs at error level
func Errorf(s string, v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(fmt.Sprintf(s, v...))
}

// Fatalf output logs at fatal level
func Fatalf(s string, v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalf(fmt.Sprintf(s, v...))
}

// Panicf output logs at fatal level
func Panicf(s string, v ...interface{}) {
	setPrefix(PANIC)
	logger.Panicf(s, v...)
}

// setPrefix set the prefix of the log output
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][ %s:%d ]", levelColors[level](levelFlags[level]), filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelColors[level](levelFlags[level]))
	}

	logger.SetPrefix(logPrefix)
}
