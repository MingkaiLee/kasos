package util

import (
	"log"
	"os"
)

var (
	// Log is the logger used by the package.
	loggerInfo  *log.Logger
	loggerError *log.Logger
)

func init() {
	loggerInfo = log.New(os.Stdout, "[INFO]", log.LstdFlags|log.Lshortfile)
	loggerError = log.New(os.Stdout, "[ERROR]", log.LstdFlags|log.Lshortfile)
}

func LogInfof(format string, v ...interface{}) {
	loggerInfo.Printf(format, v...)
}

func LogErrorf(format string, v ...interface{}) {
	loggerError.Printf(format, v...)
}
