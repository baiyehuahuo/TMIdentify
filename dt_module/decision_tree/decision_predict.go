package decision_tree

// PredictFloatFeature Predict the result of float eigenvalues
func (tree *DecisionTree) PredictFloatFeature(feature []float64) (result bool) {
	return tree.PredictIntFeature(tree.ConvertFeatureOneByBorders(feature))
}

// PredictIntFeature Predict the result of int eigenvalues
func (tree *DecisionTree) PredictIntFeature(feature []int) (result bool) {
	root := tree.Tree
	for !root.IsLeaf {
		root = root.SubTrees[feature[root.Index]]
	}
	result = root.IsTencentMeeting
	return
}

// PredictAll Classify 2D features
func (tree *DecisionTree) PredictAll(features [][]int) (result []bool) {
	result = make([]bool, 0, len(features))
	for i := range features {
		result = append(result, tree.PredictIntFeature(features[i]))
	}
	return
}

// PredictAllByIndex Predict all results according to the specified index Only predict the features of the indexes
func (tree *DecisionTree) PredictAllByIndex(features [][]int, indexes []int) (result []bool) {
	result = make([]bool, len(features))
	for _, index := range indexes {
		result[index] = tree.PredictIntFeature(features[index])
	}
	return
}
