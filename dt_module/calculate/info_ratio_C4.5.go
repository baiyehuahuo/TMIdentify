package calculate

import (
	"math"
)

func CloseInfoGainRatio(indexes []int, labels []bool) func([][]int, []int, []bool) float64 {
	entD := CalcShannonEnt(labels, indexes)
	return func(splitIndexes [][]int, indexes []int, labels []bool) (gainRatio float64) {
		gainRatio = commonInfoGainRatio(splitIndexes, indexes, labels, entD)
		return gainRatio
	}
}

func commonInfoGainRatio(splitIndexes [][]int, indexes []int, labels []bool, entD float64) (gainRatio float64) {
	var prob, ent, splitInfo float64
	datasetSize := float64(len(indexes))

	for _, splitIndex := range splitIndexes {
		prob = float64(len(splitIndex)) / datasetSize
		ent += prob * CalcShannonEnt(labels, splitIndex)
		splitInfo = splitInfo - prob*math.Log2(prob+(math.SmallestNonzeroFloat64))
	}

	if splitInfo == 0 {
		gainRatio = 0
	} else {
		gainRatio = (entD - ent) / splitInfo
	}

	return gainRatio
}
