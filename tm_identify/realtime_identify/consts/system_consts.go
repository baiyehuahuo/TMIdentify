package consts

import "time"

const (
	DefaultYamlPath = "save_option.yaml"

	DeviceNameForLinux   = "ens33"
	DeviceNameForOpenWRT = "br-lan"
	DeviceNameForWindows = "\\Device\\NPF_{363B2123-859F-4299-8552-30C5E788B3D9}"
	SnapLen              = 1024
	Promiscuous          = false
	TimeOut              = -1 * time.Second

	SystemLogPath = "systemlogs"
	LogFilePath   = SystemLogPath + "/log.txt"

	CacheExpiration = time.Minute * 5
)
