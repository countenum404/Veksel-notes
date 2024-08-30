package logger

import (
	"io/fs"
	"log"
	"os"
	"sync"
)

const (
	INFO  = "info.log"
	WARN  = "warn.log"
	ERROR = "error.log"
)

type Logger struct {
	info     *log.Logger
	warn     *log.Logger
	err      *log.Logger
	fileInfo *os.File
	fileWarn *os.File
	fileErr  *os.File
}

var (
	logger *Logger
	once   sync.Once
)

func GetLogger() *Logger {
	once.Do(func() {
		logger = newLogger()
	})

	return logger
}

func newLogger() *Logger {
	logger := &Logger{}

	fileFlags, access := os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666
	fileInfo, _ := os.OpenFile(INFO, fileFlags, fs.FileMode(access))
	fileWarn, _ := os.OpenFile(WARN, fileFlags, fs.FileMode(access))
	fileErr, _ := os.OpenFile(ERROR, fileFlags, fs.FileMode(access))
	logger.fileInfo = fileInfo
	logger.fileWarn = fileWarn
	logger.fileErr = fileErr

	logFlags := log.LstdFlags
	logger.info = log.New(fileInfo, "[INFO]", logFlags)
	logger.warn = log.New(fileWarn, "[WARN]", logFlags)
	logger.err = log.New(fileErr, "[ERROR]\t", logFlags)

	return logger
}

func (l *Logger) Info(v ...any) {
	l.info.Println(v...)
	l.info.SetOutput(os.Stdout)
	l.info.Println(v...)
	l.info.SetOutput(l.fileInfo)
}

func (l *Logger) Warn(v ...any) {
	l.warn.Println(v...)
	l.warn.SetOutput(os.Stdout)
	l.warn.Println(v...)
	l.warn.SetOutput(l.fileInfo)
}

func (l *Logger) Err(v ...any) {
	l.warn.Println(v...)
	l.warn.SetOutput(os.Stdout)
	l.warn.Println(v...)
	l.warn.SetOutput(l.fileInfo)
}
