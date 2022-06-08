package decision_tree

import (
	"discretization/calculate"
	"discretization/config"
	"discretization/data_process"
	"log"
)

type TreeNode struct {
	Index            int
	IsLeaf           bool
	IsTencentMeeting bool
	SubTrees         []*TreeNode
}

// CreateTree Generate decision tree
func CreateTree(features [][]int, labels []bool, indexes []int, usedFeature []bool,
	closeInfoFunc func(indexes []int, labels []bool) func([][]int, []int, []bool) float64) (root *TreeNode) {
	if len(indexes) == 0 {
		log.Println("error: indexes is empty") // 不应该有
		return nil
	}

	root = &TreeNode{}
	var bestGain float64
	infoFunc := closeInfoFunc(indexes, labels)
	root.Index, bestGain = selectBestFeatureIndex(features, labels, indexes, usedFeature, infoFunc)
	if bestGain == 0 || isLeaf(features, root.Index, labels, indexes, usedFeature) {
		labelCounter := calculate.LabelCounter(indexes, labels)
		return &TreeNode{
			Index:            -1,
			IsLeaf:           true,
			IsTencentMeeting: labelCounter[1] >= labelCounter[0],
			SubTrees:         nil,
		}
	}
	splitIndexes := data_process.SplitDatasetByVal(features, root.Index, indexes)
	usedFeature[root.Index] = true
	for _, splitIndex := range splitIndexes {
		root.SubTrees = append(root.SubTrees, CreateTree(features, labels, splitIndex, usedFeature, closeInfoFunc))
	}
	usedFeature[root.Index] = false
	return root
}

// selectBestFeatureIndex  Obtain optimal classification features
func selectBestFeatureIndex(features [][]int, labels []bool, indexes []int, usedFeature []bool,
	infoFunc func(splitIndexes [][]int, indexes []int, labels []bool) (gain float64)) (bestFeatureIndex int, bestFeatureGain float64) {
	bestFeatureIndex, bestFeatureGain = -1, 0.0
	var inited bool
	var featureGain float64
	var splitIndexes [][]int
	for i := range usedFeature {
		if usedFeature[i] {
			continue
		}
		splitIndexes = data_process.SplitDatasetByVal(features, i, indexes)
		featureGain = infoFunc(splitIndexes, indexes, labels)
		if featureGain > bestFeatureGain || !inited {
			bestFeatureIndex = i
			bestFeatureGain = featureGain
			inited = true
		}
	}
	return bestFeatureIndex, bestFeatureGain
}

// isLeaf judge is leaf node or not
func isLeaf(features [][]int, axis int, labels []bool, indexes []int, usedFeature []bool) (result bool) {
	labelCounter := calculate.LabelCounter(indexes, labels)
	featureCounter := calculate.FeatureCounter(features, axis, indexes)
	usedCounter := 0
	for _, used := range usedFeature {
		if used {
			usedCounter++
		} else {
			break
		}
	}
	for _, count := range labelCounter {
		result = result || count == 0
	}
	for _, count := range featureCounter {
		result = result || count == 0
	}
	result = result || usedCounter == len(usedFeature) || usedCounter > config.GlobalConfig.MaxDeeper
	return result
}
