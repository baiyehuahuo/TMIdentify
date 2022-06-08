package data_process

import "math/rand"

// CreateSequenceIndex Create a sequence index slice
func CreateSequenceIndex(num int) (sequenceIndex []int) {
	sequenceIndex = make([]int, num)
	for i := range sequenceIndex {
		sequenceIndex[i] = i
	}
	return sequenceIndex
}

func shuffle(slice []int) {
	var t int
	for i := len(slice) - 1; i > 0; i-- {
		t = rand.Intn(i)
		slice[i], slice[t] = slice[t], slice[i]
	}
}

//
func OriginTrainTestIndexCreate(dataNum int, testRate int) (trainIndexes, testIndexes []int) {
	indexes := CreateSequenceIndex(dataNum)
	testNum := dataNum / testRate
	shuffle(indexes)
	trainIndexes = indexes[:dataNum-testNum]
	testIndexes = indexes[dataNum-testNum:]
	return trainIndexes, testIndexes
}

func CrossValidationIndexCreate(dataNum, k int) ([][]int, [][]int) {
	originIndexes := CreateSequenceIndex(dataNum)
	shuffle(originIndexes)
	trainIndexes := make([][]int, 0, k)
	testIndexes := make([][]int, 0, k)
	length := len(originIndexes)
	trainLength, testLength := length-length/k, length/k
	var train, test []int
	for i := 0; i < k; i++ {
		train = make([]int, trainLength)
		test = make([]int, testLength)
		copy(train, originIndexes[:i*testLength])
		copy(test, originIndexes[i*testLength:(i+1)*testLength])
		copy(train[i*testLength:], originIndexes[(i+1)*testLength:])
		trainIndexes = append(trainIndexes, train)
		testIndexes = append(testIndexes, test)
	}
	return trainIndexes, testIndexes
}
