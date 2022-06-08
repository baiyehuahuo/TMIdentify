package calculate

import (
	"discretization/config"
	"fmt"
)

// GetScores 获得点数
func GetScores(labels []bool, predict []bool, indexes []int) (scores [4]float64) {
	TP, FN, FP, TN := GetConfusionMatrix(labels, predict, indexes)
	if config.GlobalConfig.PrintMatrix {
		fmt.Println("TP FN FP TN: ", TP, FN, FP, TN)
	}

	accuracy := float64(TP+TN) / float64(TP+TN+FP+FN)
	precision := float64(TP) / float64(TP+FP)
	recall := float64(TP) / float64(TP+FN)
	f1 := (2 * precision * recall) / (precision + recall)
	scores = [4]float64{accuracy, precision, recall, f1}
	return scores
}

func GetConfusionMatrix(labels []bool, predict []bool, indexes []int) (TP, FN, FP, TN int) {
	for _, index := range indexes {
		switch {
		case labels[index] && predict[index]:
			TP++
		case labels[index] && !predict[index]:
			FN++
		case !labels[index] && predict[index]:
			FP++
		case !labels[index] && !predict[index]:
			TN++
		}
	}
	return
}
