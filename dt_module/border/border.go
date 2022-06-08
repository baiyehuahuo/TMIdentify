package border

import (
	"discretization/config"
	"discretization/csv_read"
	"discretization/data_process"
	"fmt"
	"sort"
)

// GetBordersByInformationGain 通过计算信息熵，得到最好的边界
func GetBordersByInformationGain(csvFilePath string, closeInfoFunc func(indexes []int, labels []bool) func([][]int, []int, []bool) float64) (bestBorders []float64) {
	columnNum := csv_read.GetCSVColumnNum(csvFilePath)
	bestBorders = make([]float64, 0, columnNum-2)
	features, labels := csv_read.ReadCSVAll(csvFilePath)
	var splitIndexes [][]int
	var borders []float64
	var bestBorder, bestBorderGain float64
	var tBorderGain float64
	var i, j int
	indexes := make([]int, len(labels))
	for i := range indexes {
		indexes[i] = i
	}
	infoFunc := closeInfoFunc(indexes, labels)
	for i = 0; i < columnNum-2; i++ {
		borders = getBorders(features, i)
		if len(borders) == 0 {
			bestBorders = append(bestBorders, -1)
			continue
		}
		splitIndexes = data_process.SplitDatasetByBorder(features, i, indexes, borders[0])
		bestBorder, bestBorderGain = borders[0], infoFunc(splitIndexes, indexes, labels)
		for j = range borders {
			splitIndexes = data_process.SplitDatasetByBorder(features, i, indexes, borders[j])
			tBorderGain = infoFunc(splitIndexes, indexes, labels)
			if tBorderGain > bestBorderGain {
				bestBorder, bestBorderGain = borders[j], tBorderGain
			}
		}
		bestBorders = append(bestBorders, bestBorder)
		if config.GlobalConfig.PrintBestGain {
			fmt.Println("best gain: ", i, bestBorder, bestBorderGain)
		}
	}
	return bestBorders
}

// getBorders 获取所有可能的边界 从小到大排列
func getBorders(features [][]float64, axis int) (borders []float64) {
	set := make(map[float64]struct{})
	for i := range features {
		set[features[i][axis]] = struct{}{}
	}
	length := len(set)
	vals := make([]float64, 0, length)
	for val := range set {
		vals = append(vals, val)
	}
	sort.Float64s(vals)
	borders = make([]float64, 0, length-1)
	for i := 1; i < length; i++ {
		borders = append(borders, (vals[i]+vals[i-1])/2)
	}
	return borders
}
