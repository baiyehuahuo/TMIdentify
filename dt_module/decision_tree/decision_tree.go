package decision_tree

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

type DecisionTree struct {
	Borders   []float64
	Tree      *TreeNode
	Accuracy  float64
	Precision float64
	Recall    float64
	F1        float64
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

// ConvertFeatureByBordersAll Convert all features to int
func (tree *DecisionTree) ConvertFeatureByBordersAll(originFeatures [][]float64) (features [][]int) {
	features = make([][]int, len(originFeatures))
	for i := range originFeatures {
		features[i] = tree.ConvertFeatureOneByBorders(originFeatures[i])
	}
	return features
}

// Save decision tree
func (tree *DecisionTree) Save(savePath string) {
	file, err := os.Create(savePath)
	if err != nil {
		fmt.Println(err)
	}
	enc := gob.NewEncoder(file)
	if err = enc.Encode(tree); err != nil {
		fmt.Println(err)
	}
}

// LoadDecisionTree Load decision tree
func LoadDecisionTree(loadPath string) (tree *DecisionTree) {
	file, err := os.Open(loadPath)
	if err != nil {
		log.Fatal(err)
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&tree)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return tree
}
