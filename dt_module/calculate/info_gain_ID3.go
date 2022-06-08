package calculate

import (
	"math"
)

func CalcShannonEnt(labels []bool, indexes []int) float64 {
	labelCounter := make([]int, 2)
	for _, index := range indexes {
		if labels[index] {
			labelCounter[1]++
		} else {
			labelCounter[0]++
		}
	}
	length := float64(len(indexes))
	if length == 0 {
		return 0
	}
	var prob, shannonEnt float64
	for i := range labelCounter {
		prob = float64(labelCounter[i]) / length
		shannonEnt -= prob * math.Log2(prob+(math.SmallestNonzeroFloat64)) // 防止为0
	}
	return shannonEnt
}

// Close 是闭包计算 为了减少函数ent的计算

// CloseInfoGain 带闭包参数信息增益函数
func CloseInfoGain(indexes []int, labels []bool) func([][]int, []int, []bool) float64 {
	entD := CalcShannonEnt(labels, indexes)
	return func(splitIndexes [][]int, indexes []int, labels []bool) (gain float64) {
		gain = commonInfoGain(splitIndexes, indexes, labels, entD)
		return gain
	}
}

func commonInfoGain(splitIndexes [][]int, indexes []int, labels []bool, entD float64) (gain float64) {
	var prob, ent float64
	datasetSize := float64(len(indexes))

	for _, splitIndex := range splitIndexes {
		prob = float64(len(splitIndex)) / datasetSize
		ent += prob * CalcShannonEnt(labels, splitIndex)
	}
	gain = entD - ent
	return gain
}
