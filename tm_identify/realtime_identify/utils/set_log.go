package utils

import (
	"identifyTencentMeeting/consts"
	"log"
	"os"
)

// SetLog 设置log日志路径
func SetLog() error {
	logFile, err := os.OpenFile(consts.LogFilePath, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return nil
}
