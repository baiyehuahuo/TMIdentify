package main

import (
	"flag"
	"fmt"
	"identifyTencentMeeting/config"
	"identifyTencentMeeting/utils"
	"os"
	"strings"
)

func main() {
	// getMLScores()
	flag.IntVar(&config.GlobalConfig.MLJudgePacket, "ml", config.GlobalConfig.MLJudgePacket, "")
	// flag.StringVar(&config.GlobalConfig.DirPath, "d", config.GlobalConfig.DirPath, "")
	flag.Parse()
	config.UpdatePath()
	Initialize()
	TP, FN := readFromPcapDirectoryIntegrate(config.GlobalConfig.TMPath, true)
	FP, TN := readFromPcapDirectoryIntegrate(config.GlobalConfig.OtherPath, false)
	Destroy()
	utils.GetScores(TP, FN, FP, TN, true)
}

func getMLScores() {
	gConfig := config.GlobalConfig
	start, end := 3, 25
	var scores = make([]string, 0, end-start+1)
	for i := start; i <= end; i++ {
		fmt.Println(i)
		gConfig.MLJudgePacket = i
		config.UpdatePath()
		Initialize()
		TP, FN := readFromPcapDirectoryIntegrate(config.GlobalConfig.TMPath, true)
		FP, TN := readFromPcapDirectoryIntegrate(config.GlobalConfig.OtherPath, false)
		Destroy()
		scores = append(scores, utils.GetScores(TP, FN, FP, TN, false))
	}
	// fmt.Println(strings.Join(scores, "\n"))
	os.Exit(0)
}

func getDPIScores() {
	gConfig := config.GlobalConfig
	for precision := 0.1; precision <= 0.5; precision += 0.1 {
		gConfig.DPIFilterPrecision = precision
		for threshold := 0.1; threshold <= 0.5; threshold += 0.1 {
			fmt.Printf("\nprecision: %.1f\t threshold: %.1f\n", precision, threshold)
			gConfig.DPIBaseThreshold = threshold
			var scores = make([]string, 0, 5)
			for i := 1; i <= 5; i++ {
				gConfig.DPIJudgePacket = i
				config.UpdatePath()
				Initialize()
				TP, FN := readFromPcapDirectoryIntegrate(config.GlobalConfig.TMPath, true)
				FP, TN := readFromPcapDirectoryIntegrate(config.GlobalConfig.OtherPath, false)
				Destroy()
				// fmt.Println(TP, FN, FP, TN)
				scores = append(scores, utils.GetScores(TP, FN, FP, TN, false))
			}
			fmt.Println(strings.Join(scores, "\n"))
		}
	}

	os.Exit(0)
}

/*
go run . -ml 2
go run . -ml 3
go run . -ml 4
go run . -ml 5
go run . -ml 6
go run . -ml 7
go run . -ml 8
go run . -ml 9
go run . -ml 10
go run . -ml 11
go run . -ml 12
go run . -ml 13
go run . -ml 14
go run . -ml 15


go run . -ml 16
go run . -ml 17
go run . -ml 18
go run . -ml 19
go run . -ml 20
go run . -ml 21
go run . -ml 22
go run . -ml 23
go run . -ml 24
go run . -ml 25
*/
