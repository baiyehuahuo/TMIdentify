package main

import (
	"discretization/border"
	"discretization/calculate"
	"discretization/config"
	"discretization/csv_read"
	"discretization/data_process"
	"discretization/decision_tree"
	"fmt"
	"strings"
)

func TrainDecisionTree(mode int, csvPath string, closeFunc func([]int, []bool) func([][]int, []int, []bool) float64) {
	gConfig := config.GlobalConfig

	var decisionTree decision_tree.DecisionTree
	if !gConfig.UseInputBorders {
		decisionTree.Borders = border.GetBordersByInformationGain(csvPath, closeFunc)
	} else {
		decisionTree.Borders = gConfig.Borders[mode] //  节省border训练时间
	}
	if gConfig.PrintCalBorder {
		printSlice(decisionTree.Borders)
	}
	originFeatures, labels := csv_read.ReadCSVAll(csvPath)
	features := decisionTree.ConvertFeatureByBordersAll(originFeatures)
	trainIndexes, testIndexes := data_process.OriginTrainTestIndexCreate(len(features), 10)
	usedFeature := make([]bool, len(features[0]))
	calculate.FeatureFilter(features, trainIndexes, usedFeature)
	decisionTree.Tree = decision_tree.CreateTree(features, labels, trainIndexes, usedFeature, closeFunc)

	predictTest := decisionTree.PredictAllByIndex(features, testIndexes)
	scores := calculate.GetScores(labels, predictTest, testIndexes)
	decisionTree.Accuracy, decisionTree.Precision, decisionTree.Recall, decisionTree.F1 = scores[0], scores[1], scores[2], scores[3]
	fmt.Println(decisionTree.Accuracy, decisionTree.Precision, decisionTree.Recall, decisionTree.F1)

	if gConfig.SaveModel {
		decisionTree.Save(gConfig.SavePath[mode])
		fmt.Println(gConfig.SavePath[mode], "saved.")
	}
}

func printSlice(slice interface{}) {
	str := fmt.Sprintln(slice)
	length := len(str)
	str = str[:length-1] + ",\n"
	fmt.Printf(strings.Join(strings.Split(str, " "), ", ")) // 保存borders
}
