package calculate

import "discretization/config"

func FeatureCounter(features [][]int, axis int, indexes []int) (counter []int) {
	counter = make([]int, 2)
	for _, index := range indexes {
		counter[features[index][axis]]++
	}
	return counter
}

func LabelCounter(indexes []int, labels []bool) (counter []int) {
	counter = make([]int, 2)
	for _, index := range indexes {
		if !labels[index] {
			counter[0]++
		} else {
			counter[1]++
		}
	}
	return counter
}

func FeatureFilter(features [][]int, indexes []int, usedFeature []bool) {
	var counter []int
	for i := range usedFeature {
		counter = FeatureCounter(features, i, indexes)
		for _, count := range counter {
			if count == 0 {
				usedFeature[i] = true
				config.GlobalConfig.MaxDeeper++
				break
			}
		}
	}
}
