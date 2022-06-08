package utils

import (
	"fmt"
	"math"
)

func GetScores(TP, FN, FP, TN int, printMessage bool) string {
	if printMessage {
		fmt.Printf("TP FN:\t%d, %d\t%d\n", TP, FN, TP+FN)
		fmt.Printf("FP TN:\t%d, %d\t%d\n", FP, TN, FP+TN)
		fmt.Println("total:\t", TP+FN+FP+TN)
	}
	accuracy := float64(TP+TN) / float64(TP+TN+FP+FN)
	precision := float64(TP) / float64(TP+FP)
	recall := float64(TP) / float64(TP+FN)
	f1 := (2 * precision * recall) / (precision + recall)
	if printMessage {
		fmt.Printf("Accuracy:\t%f\n", accuracy)
		fmt.Printf("Precision:\t%f\n", precision)
		fmt.Printf("Recall:\t\t%f\n", recall)
		fmt.Printf("F1:\t\t%f\n", f1)
	}
	if math.IsNaN(precision) {
		precision = 0
	}
	if math.IsNaN(f1) {
		f1 = 0
	}
	return fmt.Sprintf("\t\t\t[%.6f,%.6f,%.6f,%.6f],", accuracy, precision, recall, f1)
}
