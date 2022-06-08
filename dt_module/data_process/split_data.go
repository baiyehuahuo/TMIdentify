package data_process

// splitDataset 按照oldIndex 和 border 分割出新的indexes
func SplitDatasetByBorder(features [][]float64, axis int, oldIndex []int, border float64) (splitIndexes [][]int) {
	splitIndexes = make([][]int, 2)
	for i := range oldIndex {
		if features[oldIndex[i]][axis] < border {
			splitIndexes[0] = append(splitIndexes[0], oldIndex[i])
		} else {
			splitIndexes[1] = append(splitIndexes[1], oldIndex[i])
		}
	}
	return
}

// SplitDatasetByVal 按照 0 1 划分
func SplitDatasetByVal(features [][]int, axis int, oldIndex []int) (splitIndexes [][]int) {
	splitIndexes = make([][]int, 2)
	for _, index := range oldIndex {
		splitIndexes[features[index][axis]] = append(splitIndexes[features[index][axis]], index)
	}
	return
}
