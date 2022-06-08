package main

import (
	"dpi_module/config"
	"dpi_module/types"
	"dpi_module/utils"
	"flag"
	"fmt"
)

func main() {
	gConfig := config.GlobalConfig
	flag.IntVar(&gConfig.Mode, "m", gConfig.Mode, "")
	flag.StringVar(&gConfig.CsvPath, "p", gConfig.CsvPath, "")
	flag.StringVar(&gConfig.DirPath, "d", gConfig.DirPath, "")
	flag.Float64Var(&gConfig.BaseThreshold, "base", gConfig.BaseThreshold, "")
	flag.IntVar(&gConfig.JudgeThreshold, "judge", gConfig.JudgeThreshold, "")
	flag.Parse()
	extraSig := "[46 113 113 46 99 111 109 48 14 6 3 85 29 15 1 1 255 4 4 3 2 5 160 48 29 6 3 85 29 37 4 22 48 20 6 8 43 6 1 5 5 7 3 1 6 8 43 6 1 5 5 7 3 2 48 68 6 3 85 29 31 4 61 48 59\n 48 57 160 55 160 53 134 51 104 116 116 112 58 47 47 99 114 108 46 100 105 103 105 99 101 114 116 46 99 110 47 68 105 103 105 67 101 114 116 83 101 99 117 114 101 83 105 116 101 67 78 67 65 71 51 46 99 114 108 48 62 6 3 85 29 32 4 55 48 53 48 51 \n6 6 103 129 12 1 2 2 48 41 48 39 6 8 43 6 1 5 5 7 2 1 22 27 104 116 116 112 58 47 47 119 119 119 46 100 105 103 105 99 101 114 116 46 99 111 109 47 67 80 83 48 120 6 8 43 6 1 5 5 7 1 1 4 108 48 106 48 35 6 8 43 6 1 5 5 7 48 1 134 23 104 116 116 1\n12 58 47 47 111 99 115 112 46 100 105 103 105 99 101 114 116 46 99 110 48 67 6 8 43 6 1 5 5 7 48 2 134 55 104 116 116 112 58 47 47 99 97 99 101 114 116 115 46 100 105 103 105 99 101 114 116 46 99 110 47 68 105 103 105 67 101 114 116 83 101 99 117\n 114 101 83 105 116 101 67 78 67 65 71 51 46 99 114 116 48 12 6 3 85 29 19 1 1 255 4 2 48 0 48 130 1]"
	switch gConfig.Mode {
	case 0:
		utils.ReadPcapFile(gConfig.CsvPath)
		cspExtract()
	case 1:
		utils.ReadPcapDir(gConfig.DirPath)
		cspExtract()
	case 2:
		calculatePrecision(utils.LoadSigs(utils.GetRawSaveName()), gConfig.DirPath, gConfig.OtherPath)
	case 3:
		utils.ReadPcapDir(gConfig.DirPath)
		cspExtract()
		calculatePrecision(utils.LoadSigs(utils.GetRawSaveName()), gConfig.DirPath, gConfig.OtherPath)
	case 4:
		totalEvaluate(utils.LoadSigs(utils.GetFilterSaveName()), gConfig.DirPath, gConfig.OtherPath)
	case 5:
		utils.PrintSigBytes(utils.LoadSigs(utils.GetFilterSaveName()))
		fmt.Println(utils.GetFilterSaveName())
	case 6:
		deleteSig(utils.LoadSigs(utils.GetFilterSaveName()), utils.PrintSliceSigToString(extraSig))
	case 7:
		addSig(utils.LoadSigs(utils.GetFilterSaveName()), utils.PrintSliceSigToString(extraSig))
	default:
		test2()
	}
}

func test() {
	gConfig := config.GlobalConfig
	utils.ReadPcapDir(gConfig.DirPath)
	// flowMap := types.GetFlowMap()
	// fmt.Println(len(flowMap))

	payloads := types.GetFlowPayloads()
	countMap := make(map[int]int)
	for _, payload := range payloads {
		countMap[len(payload)]++
	}
	utils.PrintSplitCountMap(countMap, "total")
	fmt.Println(countMap)
}

func test2() {
	gConfig := config.GlobalConfig
	for minPrecision := 0.1; minPrecision <= 0.5; minPrecision += 0.1 {
		gConfig.MinPrecision = minPrecision
		for judgeThreshold := 1; judgeThreshold <= 5; judgeThreshold++ {
			gConfig.JudgeThreshold = judgeThreshold
			for baseThreshold := 0.1; baseThreshold <= 0.5; baseThreshold += 0.1 {
				gConfig.BaseThreshold = baseThreshold
				// utils.ReadPcapDir(gConfig.DirPath)
				// cspExtract()
				// types.Clean()
				calculatePrecision(utils.LoadSigs(utils.GetRawSaveName()), gConfig.DirPath, gConfig.OtherPath)
				fmt.Println(utils.GetFilterSaveName())
			}
		}
	}
}

/*
go run . -m 3 -base 0.2 -judge 1
go run . -m 3 -base 0.3 -judge 1
go run . -m 3 -base 0.4 -judge 1
go run . -m 3 -base 0.5 -judge 1

*/
