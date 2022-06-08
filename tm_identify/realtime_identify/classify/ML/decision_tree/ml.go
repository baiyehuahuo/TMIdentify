package decision_tree

import (
	"encoding/gob"
	"os"
)

type TreeNode struct {
	Index            int
	IsLeaf           bool
	IsTencentMeeting bool
	SubTrees         []*TreeNode
}

type DecisionTree struct {
	Borders   []float64
	Tree      *TreeNode
	Accuracy  float64
	Precision float64
	Recall    float64
	F1        float64
}

type RandomForest []*DecisionTree

var decisionTrees []*DecisionTree
var rf RandomForest

// LoadDecisionTree Load decision tree
func LoadDecisionTree(loadPath string) (tree *DecisionTree, err error) {
	file, err := os.Open(loadPath)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&tree)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func LoadMLModel(mlPaths []string) error {
	for i := range mlPaths {
		if tree, err := LoadDecisionTree(mlPaths[i]); err == nil {
			decisionTrees = append(decisionTrees, tree)
		} else {
			return err
		}
	}
	rf = RandomForest(decisionTrees)
	return nil
}

func Classify(features []float64, mode string) bool {
	switch mode {
	case "ID3":
		return decisionTrees[0].PredictFloatFeature(features)
	case "C4.5":
		return decisionTrees[1].PredictFloatFeature(features)
	case "CART":
		return decisionTrees[2].PredictFloatFeature(features)
	case "RF":
		count := 0
		for i := range decisionTrees {
			if decisionTrees[i].PredictFloatFeature(features) {
				count++
			}
		}
		return count >= 2
	}
	return false
}

// ConvertFeatureOneByBorders Convert single row feature to int
func (tree *DecisionTree) ConvertFeatureOneByBorders(originFeature []float64) (feature []int) {
	feature = make([]int, len(originFeature))
	for i := range originFeature {
		if originFeature[i] < tree.Borders[i] {
			feature[i] = 0
		} else {
			feature[i] = 1
		}
	}
	return feature
}

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
