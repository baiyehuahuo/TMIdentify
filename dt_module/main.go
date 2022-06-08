package main

import (
	"discretization/calculate"
	"discretization/config"
	"discretization/decision_tree"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var closeFunc = []func([]int, []bool) func([][]int, []int, []bool) float64{calculate.CloseInfoGain, calculate.CloseInfoGainRatio, calculate.CloseGini}

func main() {
	rand.Seed(time.Now().Unix())
	packetNum := -1
	flag.IntVar(&packetNum, "p", -1, "")
	flag.IntVar(&config.GlobalConfig.Mode, "m", config.GlobalConfig.Mode, "train mode.")
	flag.StringVar(&config.GlobalConfig.CSVFilePath, "c", config.GlobalConfig.CSVFilePath, "train csv source.")
	flag.Parse()

	gConfig := config.GlobalConfig
	if packetNum != -1 {
		gConfig.CSVFilePath = fmt.Sprintf("./test/NonVPN_all_%d_with_TencentMeeting.csv", packetNum)
	}
	gConfig.SavePath = config.CreateModelPath(gConfig.CSVFilePath)
	gConfig.ModelPath = config.CreateModelPath(gConfig.CSVFilePath)

	modelPaths := gConfig.ModelPath
	mode := gConfig.Mode
	switch mode {
	case 0, 1, 2:
		decision_tree.TestTree(modelPaths[mode])
	case 3:
		decision_tree.TestForest(modelPaths)
	case 4, 5, 6:
		mode -= 4
		decision_tree.GetBetterDecision(modelPaths[mode], gConfig.TestRate, closeFunc[mode])
	case 7, 8, 9:
		mode -= 7
		TrainDecisionTree(mode, gConfig.CSVFilePath, closeFunc[mode])
	case 10:
		for i := range closeFunc {
			decision_tree.GetBetterDecision(modelPaths[i], gConfig.TestRate, closeFunc[i])
		}
	case 11:
		for i := range closeFunc {
			TrainDecisionTree(i, gConfig.CSVFilePath, closeFunc[i])
		}
	case 12:
		for i := range closeFunc {
			decision_tree.CrossValidation(gConfig.CSVFilePath, gConfig.TestRate, closeFunc[i])
		}
	default:
		test()
	}
}

func test() {
	gConfig := config.GlobalConfig
	strs := []string{"ID3", "C4.5", "CART"}
	var modelPaths []string
	for j := range closeFunc {
		fmt.Println(strs[j])
		for i := 3; i <= 25; i++ {
			gConfig.CSVFilePath = fmt.Sprintf("./test/NonVPN_all_%d_with_TencentMeeting.csv", i)
			// fmt.Println(gConfig.CSVFilePath)
			gConfig.SavePath = config.CreateModelPath(gConfig.CSVFilePath)
			gConfig.ModelPath = config.CreateModelPath(gConfig.CSVFilePath)
			modelPaths = gConfig.ModelPath
			// TrainDecisionTree(j, gConfig.CSVFilePath, closeFunc[j])
			// decision_tree.GetBetterDecision(modelPaths[j], gConfig.TestRate, closeFunc[j])
			decision_tree.TestTree(modelPaths[j])
			// decision_tree.TestForest(modelPaths)
		}
		fmt.Printf("\n\n\n\n\n\n")
	}

	fmt.Println("RF")
	for i := 3; i <= 25; i++ {
		gConfig.CSVFilePath = fmt.Sprintf("./test/NonVPN_all_%d_with_TencentMeeting.csv", i)
		// fmt.Println(gConfig.CSVFilePath)
		gConfig.SavePath = config.CreateModelPath(gConfig.CSVFilePath)
		gConfig.ModelPath = config.CreateModelPath(gConfig.CSVFilePath)
		modelPaths = gConfig.ModelPath
		decision_tree.TestForest(modelPaths)
	}
}

/*
go run . -m 11 -c test/NonVPN_all_5_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_5_with_TencentMeeting.csv

go run . -m 11 -c test/NonVPN_all_5_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_5_with_TencentMeeting.csv

go run . -m 11 -c test/NonVPN_all_5_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_5_with_TencentMeeting.csv

go run . -m 11 -c test/NonVPN_all_5_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_5_with_TencentMeeting.csv

go run . -m 11 -c test/NonVPN_all_5_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_5_with_TencentMeeting.csv

go run . -m 11 -c test/NonVPN_all_5_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_5_with_TencentMeeting.csv


go run . -m 11 -c test/NonVPN_all_5_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_5_with_TencentMeeting.csv

go run . -m 11 -c test/NonVPN_all_10_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_10_with_TencentMeeting.csv

go run . -m 11 -c test/NonVPN_all_15_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_15_with_TencentMeeting.csv

go run . -m 11 -c test/NonVPN_all_17_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_17_with_TencentMeeting.csv


go run . -m 11 -c test/NonVPN_all_20_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_20_with_TencentMeeting.csv

go run . -m 11 -c test/NonVPN_all_25_with_TencentMeeting.csv
go run . -m 10 -c test/NonVPN_all_25_with_TencentMeeting.csv

*/
