package calculate

// 返回负数是为了能用更大的gini系数进行分类

// CloseGini 带闭包参数信息增益函数
func CloseGini(indexes []int, labels []bool) func([][]int, []int, []bool) float64 {
	return func(splitIndexes [][]int, indexes []int, labels []bool) (giniD float64) {
		giniD = commonGini(splitIndexes, indexes, labels)
		return giniD
	}
}

func commonGini(splitIndexes [][]int, indexes []int, labels []bool) (giniD float64) {
	var prob, gini float64
	datasetSize := float64(len(indexes))

	var counter []int
	var length float64
	for _, splitIndex := range splitIndexes {
		if length = float64(len(splitIndex)); length == 0 {
			continue
		}
		counter = LabelCounter(splitIndex, labels)
		gini = 1
		for _, count := range counter {
			prob = float64(count) / length
			gini -= prob * prob
		}
		giniD += (length / datasetSize) * gini
	}
	giniD = -giniD
	return giniD
}
