package consts

import "time"

const (
	DefaultCSVLabel = "No Label"
	DataPath        = "./data"
	DefaultYamlPath = "save_option.yaml"

	SystemLogPath = "systemlogs"
	LogFilePath   = SystemLogPath + "/log.txt"

	CacheExpiration = time.Minute * 5
)
