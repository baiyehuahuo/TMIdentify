package csv_read

import (
	"discretization/consts"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// getCSVColumnNum 获取CSV列数
func GetCSVColumnNum(csvPath string) int {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	r := csv.NewReader(csvFile)
	header, _ := r.Read()
	return len(header)
}

// readCSVAll 获取csv所有数据转换为float64二维数组返回，标签int返回
func ReadCSVAll(csvPath string) (features [][]float64, labels []bool) {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	r := csv.NewReader(csvFile)
	if _, err = r.Read(); err != nil {
		log.Fatal(err)
	}
	var csvLine []string
	var feature []float64
	var f float64
	var label bool
	dimension := GetCSVColumnNum(csvPath) - consts.SkipHeadColumn - 1
	for {
		csvLine, err = r.Read()
		if err != nil {
			break
		}
		feature = make([]float64, 0, dimension)
		for _, data := range csvLine[consts.SkipHeadColumn : consts.SkipHeadColumn+dimension] {
			f, _ = strconv.ParseFloat(data, 64)
			feature = append(feature, f)
		}
		features = append(features, feature)
		label, _ = strconv.ParseBool(csvLine[consts.SkipHeadColumn+dimension])
		labels = append(labels, label)
	}
	return
}

// readCSVColumn 获取csv某一列数据转换为float64数组返回
func readCSVColumn(csvPath string, i int) []float64 {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	r := csv.NewReader(csvFile)
	var csvLine []string
	var result []float64
	var f float64
	_, _ = r.Read()
	for {
		csvLine, err = r.Read()
		if err != nil {
			break
		}
		if f, err = strconv.ParseFloat(csvLine[i], 64); err != nil {
			log.Fatal(err)
		}
		result = append(result, f)
	}
	return result
}

// readCSVLabel 读取csv文件标签
func readCSVLabel(csvPath string, i int) []int {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	r := csv.NewReader(csvFile)
	var csvLine []string
	var result []int
	var label int64
	_, _ = r.Read()
	for {
		csvLine, err = r.Read()
		if err != nil {
			break
		}
		if label, err = strconv.ParseInt(csvLine[i], 10, 64); err != nil {
			log.Fatal(err)
		}
		result = append(result, int(label))
	}
	return result
}
