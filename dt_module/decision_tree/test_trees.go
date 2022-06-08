// Package decision_tree
// this file is various test functions
package decision_tree

import (
	"discretization/border"
	"discretization/calculate"
	"discretization/config"
	"discretization/csv_read"
	"discretization/data_process"
	"fmt"
)

// GetBetterDecision Obtain the optimal general tree by calculating the average optimal F1 value
func GetBetterDecision(modelPath string, testRate int, closeInfoFunc func(indexes []int, labels []bool) func([][]int, []int, []bool) float64) {
	decisionTree := LoadDecisionTree(modelPath)
	originFeatures, labels := csv_read.ReadCSVAll(config.GlobalConfig.CSVFilePath)
	features := decisionTree.ConvertFeatureByBordersAll(originFeatures)
	var scores, tScores [4]float64
	for i := 0; i < config.GlobalConfig.GetBetterRebuildTimes; i++ {
		trainIndexes, testIndexes := data_process.OriginTrainTestIndexCreate(len(features), testRate)
		usedFeature := make([]bool, len(features[0]))
		decisionTree.Tree = CreateTree(features, labels, trainIndexes, usedFeature, closeInfoFunc)

		predictTest := decisionTree.PredictAllByIndex(features, testIndexes)
		scores = calculate.GetScores(labels, predictTest, testIndexes)
		for j := 1; j < config.GlobalConfig.TestModelTimes; j++ {
			_, testIndexes = data_process.OriginTrainTestIndexCreate(len(features), testRate)
			predictTest = decisionTree.PredictAllByIndex(features, testIndexes)
			tScores = calculate.GetScores(labels, predictTest, testIndexes)
			for k := range scores {
				scores[k] += tScores[k]
			}
		}

		for j := range scores {
			scores[j] /= float64(config.GlobalConfig.TestModelTimes)
		}

		if scores[3] > decisionTree.F1 {
			decisionTree.Accuracy, decisionTree.Precision, decisionTree.Recall, decisionTree.F1 = scores[0], scores[1], scores[2], scores[3]
			decisionTree.Save(modelPath)
		}
	}
	fmt.Println("max F1: ", decisionTree.Accuracy, decisionTree.Precision, decisionTree.Recall, decisionTree.F1)
}

// TestTree Test the performance of general trees
func TestTree(modelPath string) {
	decisionTree := LoadDecisionTree(modelPath)
	originFeatures, labels := csv_read.ReadCSVAll(config.GlobalConfig.CSVFilePath)
	var features = decisionTree.ConvertFeatureByBordersAll(originFeatures)
	var scores [4]float64
	var totalScores [4]float64
	for i := 0; i < config.GlobalConfig.TestModelTimes; i++ {
		_, testIndexes := data_process.OriginTrainTestIndexCreate(len(features), config.GlobalConfig.TestRate)
		predictTest := decisionTree.PredictAllByIndex(features, testIndexes)
		scores = calculate.GetScores(labels, predictTest, testIndexes)
		for i := range totalScores {
			totalScores[i] += scores[i]
		}
	}
	f := float64(config.GlobalConfig.TestModelTimes)
	fmt.Printf("[%f,%f,%f,%f],\n", totalScores[0]/f, totalScores[1]/f, totalScores[2]/f, totalScores[3]/f)
}

// TestForest Test the performance of random forests
func TestForest(modelPaths []string) {
	rf := LoadRandomForest(modelPaths)
	originFeatures, labels := csv_read.ReadCSVAll(config.GlobalConfig.CSVFilePath)
	var scores [4]float64
	var totalScore [4]float64
	for i := 0; i < config.GlobalConfig.TestModelTimes; i++ {
		_, testIndexes := data_process.OriginTrainTestIndexCreate(len(originFeatures), 10)
		predictTest := rf.PredictAllFloatByIndex(originFeatures, testIndexes)
		scores = calculate.GetScores(labels, predictTest, testIndexes)
		if config.GlobalConfig.PrintMatrix {
			fmt.Println(scores)
		}
		for j := range scores {
			totalScore[j] += scores[j]
		}
	}

	f := float64(config.GlobalConfig.TestModelTimes)
	fmt.Printf("[%f,%f,%f,%f],\n", totalScore[0]/f, totalScore[1]/f, totalScore[2]/f, totalScore[3]/f)
	// fmt.Printf("[")
	// for i := range scores {
	// 	fmt.Printf("%f", totalScore[i]/float64(config.GlobalConfig.TestModelTimes))
	// 	if i != len(scores)-1 {
	// 		fmt.Printf(",")
	// 	}
	// }
	// fmt.Println("]")
}

// CrossValidation for general trees
func CrossValidation(dataPath string, testRate int, closeFunc func([]int, []bool) func([][]int, []int, []bool) float64) {
	var decisionTree DecisionTree
	decisionTree.Borders = border.GetBordersByInformationGain(dataPath, closeFunc)
	originFeatures, labels := csv_read.ReadCSVAll(dataPath)
	features := decisionTree.ConvertFeatureByBordersAll(originFeatures)
	trainIndexesArr, testIndexesArr := data_process.CrossValidationIndexCreate(len(features), testRate)
	var totalScores [4]float64
	for i := range trainIndexesArr {
		usedFeature := make([]bool, len(features[0]))
		calculate.FeatureFilter(features, trainIndexesArr[i], usedFeature)
		decisionTree.Tree = CreateTree(features, labels, trainIndexesArr[i], usedFeature, closeFunc)
		predictTest := decisionTree.PredictAllByIndex(features, testIndexesArr[i])
		scores := calculate.GetScores(labels, predictTest, testIndexesArr[i])
		for j := range scores {
			totalScores[j] += scores[j]
		}
	}
	for i := range totalScores {
		totalScores[i] /= float64(testRate)
	}
	fmt.Println("The end of cross validation: ", totalScores)
}
