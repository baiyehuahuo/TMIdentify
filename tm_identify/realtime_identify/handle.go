package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"identifyTencentMeeting/classify/DPI"
	"identifyTencentMeeting/classify/ML/decision_tree"
	"identifyTencentMeeting/config"
	"identifyTencentMeeting/consts"
	"identifyTencentMeeting/types"
	"identifyTencentMeeting/utils"
	"log"
	"os"
	"strings"
)

// Initialize initialize all preconditions
func Initialize() {
	var err error
	if err = DirectoryCreate(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	if err = utils.SetLog(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	if err = printDevicesMessage(config.GlobalConfig.DeviceName); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Println("Your host IP is :", config.GlobalConfig.HostIP[0])
	if err = DPI.LoadSigMap(config.GlobalConfig.SigPath); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Println("Signatures loaded. ")
	// DPI.PrintSigs()
	// DPI.PrintSigs()
	if err = decision_tree.LoadMLModel(config.GlobalConfig.TreesPath); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Println("Decision Tree models loaded.")
	// fmt.Println(config.GlobalConfig.TreesPath)
	types.InitCache(consts.CacheExpiration)
}

func DirectoryCreate() error {
	if err := os.MkdirAll(consts.SystemLogPath, os.ModePerm); err != nil {
		log.Print("Create directory path fail ", err)
		return utils.ErrWrapOrWithMessage(true, err)
	}
	return nil
}

// printDevicesMessage print your need device message
func printDevicesMessage(deviceName string) error {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return utils.ErrWrapOrWithMessage(true, err)
	}
	fmt.Println("Device(s) found:")
	var ipStr string
	for _, device := range devices {
		if deviceName == device.Name {
			fmt.Println("Name: ", device.Name)
			config.GlobalConfig.HostIP = make([]string, 0, len(device.Addresses))
			for _, address := range device.Addresses {
				ipStr = address.IP.String()
				fmt.Println("- IP address: ", ipStr)
				if strings.Index(ipStr, ".") == -1 {
					config.GlobalConfig.HostIP = append(config.GlobalConfig.HostIP, ipStr)
				} else {
					config.GlobalConfig.HostIP = append([]string{ipStr}, config.GlobalConfig.HostIP...) // IPv4的地址用的比较多，作为第一条判断依据
				}
			}

			break
		}
	}
	return nil
}

func Destroy() {
	types.DestroyCache()
}
