// Package decision_tree
// this file is various functional functions of random forest
package decision_tree

type RandomForest []*DecisionTree

func LoadRandomForest(paths []string) (rf *RandomForest) {
	rf = &RandomForest{}
	for _, path := range paths {
		*rf = append(*rf, LoadDecisionTree(path))
	}
	return rf
}

// PredictFloatFeature 预测float特征值的结果
func (rf *RandomForest) PredictFloatFeature(feature []float64) (result bool) {
	isTencent := 0.0
	for _, tree := range *rf {
		if tree.PredictIntFeature(tree.ConvertFeatureOneByBorders(feature)) {
			isTencent += tree.Precision
		} else {
			isTencent -= tree.Precision
		}
	}
	result = isTencent > 0
	return result
}

// PredictAllFloatByIndex 根据指定索引预测全部结果 只对index索引的特征进行预测
func (rf *RandomForest) PredictAllFloatByIndex(features [][]float64, indexes []int) (result []bool) {
	result = make([]bool, len(features))
	for _, index := range indexes {
		result[index] = rf.PredictFloatFeature(features[index])
	}
	return result
}

// PredictIntFeature 预测int特征值的结果
func (rf *RandomForest) PredictIntFeature(feature []int) (result bool) {
	isTencent := 0.0
	for _, tree := range *rf {
		root := tree.Tree
		for !root.IsLeaf {
			root = root.SubTrees[feature[root.Index]]
		}
		if root.IsTencentMeeting {
			isTencent += tree.F1
		} else {
			isTencent -= tree.F1
		}
	}
	result = isTencent > 0
	return result
}

// PredictAll 对二维特征进行分类
func (rf *RandomForest) PredictAll(features [][]int) (result []bool) {
	result = make([]bool, 0, len(features))
	for i := range features {
		result = append(result, rf.PredictIntFeature(features[i]))
	}
	return
}

// PredictAllByIndex 根据指定索引预测全部结果 只对index索引的特征进行预测
func (rf *RandomForest) PredictAllByIndex(features [][]int, indexes []int) (result []bool) {
	result = make([]bool, len(features))
	for _, index := range indexes {
		result[index] = rf.PredictIntFeature(features[index])
	}
	return
}
