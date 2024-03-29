package util

import (
	"log"
	"os"
)

const logPath = "/mnt/log/server.txt"

var (
	// Log is the logger used by the package.
	loggerInfo  *log.Logger
	loggerError *log.Logger
)

func init() {
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	loggerInfo = log.New(logFile, "[INFO]", log.LstdFlags|log.Lshortfile)
	loggerError = log.New(logFile, "[ERROR]", log.LstdFlags|log.Lshortfile)
}

func LogInfof(format string, v ...interface{}) {
	loggerInfo.Printf(format, v...)
	loggerInfo.Println()
}

func LogErrorf(format string, v ...interface{}) {
	loggerError.Printf(format, v...)
	loggerError.Println()
}
