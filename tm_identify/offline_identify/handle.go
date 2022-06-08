package main

import (
	"encoding/csv"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"identifyTencentMeeting/classify/DPI"
	"identifyTencentMeeting/classify/ML/decision_tree"
	"identifyTencentMeeting/config"
	"identifyTencentMeeting/consts"
	"identifyTencentMeeting/types"
	"identifyTencentMeeting/utils"
	"log"
	"os"
	"path"
	"path/filepath"
)

func readFromPcapDirectoryIntegrate(dirPath string, isTM bool) (T int, F int) {
	config.GlobalConfig.Label = isTM
	if err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if path.Ext(filePath) == ".pcap" {
			// fmt.Println(filePath)
			handle, err := pcap.OpenOffline(filePath)
			if err != nil {
				log.Fatal(err)
			}
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			for packet := range packetSource.Packets() {
				types.GetFlowFromPacket(packet)
			}
			handle.Close()
			t, f := types.ResetCache()
			T, F = T+t, F+f
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return T, F
}

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
	if err = CSVFileCreate(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	if config.GlobalConfig.UseDPI {
		if err = DPI.LoadSigMap(config.GlobalConfig.SigPath); err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
	}
	if config.GlobalConfig.UseML {
		if err = decision_tree.LoadMLModel(config.GlobalConfig.TreesPath); err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
	}
	types.InitCache(consts.CacheExpiration)
}

func DirectoryCreate() error {
	if err := os.MkdirAll(consts.SystemLogPath, os.ModePerm); err != nil {
		log.Print("Create directory path fail ", err)
		return utils.ErrWrapOrWithMessage(true, err)
	}
	return nil
}

func CSVFileCreate() error {
	gConfig := config.GlobalConfig
	// gConfig.DirPath = strings.ReplaceAll(gConfig.DirPath, "\\", "/")
	// csvPath := path.Join(gConfig.DirPath, path.Base(gConfig.DirPath)) + "_" + strconv.Itoa(gConfig.MLJudgePacket) + ".csv"
	csvPath := path.Join(consts.DataPath, fmt.Sprintf("NonVPN_all_%d_with_TencentMeeting.csv", gConfig.MLJudgePacket))

	csvFile, err := os.OpenFile(csvPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return utils.ErrWrapOrWithMessage(false, err)
	}
	w := csv.NewWriter(csvFile)
	w.Comma = ',' // Just like CICFlowMeter
	w.UseCRLF = true
	header := append(append([]string{"flowKey"}, utils.GetStructFieldNames(types.Features{})...), "label")
	if err = w.Write(header); err != nil {
		return utils.ErrWrapOrWithMessage(false, err)
	}
	w.Flush()
	config.GlobalConfig.CSVFile = csvFile
	config.GlobalConfig.CSVWriter = w
	config.GlobalConfig.CSVColumn = len(header)
	return nil
}

func Destroy() {
	types.DestroyCache()
	if err := config.GlobalConfig.CSVFile.Close(); err != nil {
		log.Fatal(err)
	}
}
