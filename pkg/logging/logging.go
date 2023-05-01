package logging

import (
	"fmt"
	"os"
	"time"
)

type LOG_TYPE string

const (
	INFO    LOG_TYPE = "[INF]"
	GOOD    LOG_TYPE = "[OK ]"
	ERROR   LOG_TYPE = "[ERR]"
	WARNING LOG_TYPE = "[WRN]"
)

type Logger struct {
	LogFileHandle *os.File
}

// Return a nice human time
func (logger *Logger) getNiceTime() string {
	var nicetime = time.Now().Format("2006-01-02 15:04:05")
	return nicetime
}

// Log to file
func LogFile(logfilepath string) *Logger {
	var logfile, err = os.OpenFile(logfilepath+"jserver.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		err = os.MkdirAll(logfilepath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		logfile, err = os.Create(logfilepath + "jserver.log")
	}
	return &Logger{LogFileHandle: logfile}
}

// Log to console
func (logger *Logger) Log(str string, logType LOG_TYPE) {
	go func() {
		var log = string(logType) + " " + logger.getNiceTime() + "\t " + str
		logger.LogFileHandle.Write([]byte(log + "\n"))
		fmt.Println(log)
	}()
}