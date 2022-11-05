package logging

import (
	"fmt"
	"time"
)

type LOG_TYPE string

const (
	INFO    LOG_TYPE = "[INF]"
	GOOD    LOG_TYPE = "[OK ]"
	ERROR   LOG_TYPE = "[ERR]"
	WARNING LOG_TYPE = "[WRN]"
)

func getNiceTime() string {
	nicetime := time.Now().Format("2006-01-02 15:04:05")
	return nicetime
}
func Log(str string, logType LOG_TYPE) {
	fmt.Println(string(logType) + " " + getNiceTime() + "\t " + str)
}
