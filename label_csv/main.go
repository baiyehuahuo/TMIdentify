// 在所有流量中标记属于腾讯会议的流量
// 执行 go run . -d Test/tencentMeetingPackets.pcap_Flow.csv -t Test/AllPackets.csv
// 标记会话 根据 五元组(两个IP+两个端口+协议号) 来分类，不区分方向
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const (
	TotalPath    = "Test/AllPackets.pcap_Flow.csv"
	DicPath      = "Test/tencentMeetingPackets.pcap_Flow.csv"
	CSVPath      = "Test/AllPackets.pcap_Flow.csv"
	TempFile     = "temp.csv"
	CSVColumnNum = 84
	YamlPath     = "config.yaml"
)

func main() {
	flag.IntVar(&Config.Mode, "m", Config.Mode, "processMode")
	flag.Parse()
	// time.Sleep(time.Hour)

	switch Config.Mode {
	case 0:
		// 处理全部流与腾讯会议流的关系
		processAllByTencentMeeting(Config.DicPath, Config.TotalPath)
	case 1:
		// 为一个CSV文件打上同一个标签
		justLabelCSV(Config.CSVPath, Config.Label)
	case 2:
		// 合并csv文件
		combineArgsCSV(Config.CombineCSVPaths[:len(Config.CombineCSVPaths)-1], Config.CombineCSVPaths[len(Config.CombineCSVPaths)-1])
	case 3:
		combineDirCSV(Config.DirPath, Config.CombineCSVPath)
	case 4:
		convertFeaturesByBorders(Config.CSVPath)
	}
}

func processAllByTencentMeeting(dicPath, totalPath string) {
	tempFile := filepath.Dir(totalPath) + "/" + TempFile
	nfs, err := os.OpenFile(tempFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(nfs)

	flowDict := getFlowDict(dicPath)
	fs, err := os.Open(totalPath)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(fs)

	w.Comma = r.Comma
	w.UseCRLF = true

	columnSlice, _ := r.Read()
	if err = w.Write(columnSlice); err != nil {
		log.Fatal(err)
	}

	var ok bool
	var counter int
	var labelIndex = len(columnSlice) - 1
	for {
		columnSlice, err = r.Read()
		if err != nil {
			break
		}
		if _, ok = flowDict[columnSlice[0]]; ok {
			columnSlice[labelIndex] = "true" // 是腾讯会议流量
			counter++
		} else {
			columnSlice[labelIndex] = "false" // 不是腾讯会议流量
			if err = w.Write(columnSlice); err != nil {
				log.Fatal(err)
			}
		}

	}
	w.Flush()
	if err = fs.Close(); err != nil {
		log.Fatal(err)
	}
	if err = nfs.Close(); err != nil {
		log.Fatal(err)
	}
	if err = os.Remove(totalPath); err != nil {
		log.Fatal(err)
	}
	if err = os.Rename(tempFile, totalPath); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("success labeled %d\n", counter)
}

func justLabelCSV(csvPath string, label bool) {
	labelString := strconv.FormatBool(label)
	tempFile := filepath.Dir(csvPath) + "/" + TempFile
	nfs, err := os.OpenFile(tempFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(nfs)

	// flowDict := getFlowDict(dicPath)
	fs, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(fs)

	w.Comma = r.Comma
	w.UseCRLF = true

	columnSlice, _ := r.Read()
	if err = w.Write(columnSlice); err != nil {
		log.Fatal(err)
	}

	var counter int
	var labelIndex = len(columnSlice) - 1
	for {
		columnSlice, err = r.Read()
		if err != nil {
			break
		}
		columnSlice[labelIndex] = labelString
		counter++
		if err = w.Write(columnSlice); err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()
	if err = fs.Close(); err != nil {
		log.Fatal(err)
	}
	if err = nfs.Close(); err != nil {
		log.Fatal(err)
	}
	if err = os.Remove(csvPath); err != nil {
		log.Fatal(err)
	}
	if err = os.Rename(tempFile, csvPath); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("success labeled %d to %s\n", counter, csvPath)
}

func combineArgsCSV(csvPaths []string, savePath string) {
	nfs, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(nfs)
	w.Comma = ','
	w.UseCRLF = true
	var headWrite = true
	var counter int
	for _, csvPath := range csvPaths {
		fs, err := os.Open(csvPath)
		if err != nil {
			log.Fatal(err)
		}
		r := csv.NewReader(fs)
		columnSlice, _ := r.Read()
		if headWrite {
			if err = w.Write(columnSlice); err != nil {
				log.Fatal(err)
			}
			headWrite = false
		}
		for {
			columnSlice, err = r.Read()
			if err != nil {
				break
			}
			if err = w.Write(columnSlice); err != nil {
				log.Fatal(err)
			}
			counter++
			w.Flush()
		}
		if err = fs.Close(); err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()
	if err = nfs.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("success combine %d to %s\n", counter, savePath)
}

func combineDirCSV(dirPath, csvPath string) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	nfs, err := os.OpenFile(csvPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(nfs)
	w.Comma = ','
	w.UseCRLF = true
	var headWrite = true
	var counter int
	for _, file := range files {
		fs, err := os.Open(dirPath + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		r := csv.NewReader(fs)
		columnSlice, _ := r.Read()
		if headWrite {
			if err = w.Write(columnSlice); err != nil {
				log.Fatal(err)
			}
			headWrite = false
		}
		for {
			columnSlice, err = r.Read()
			if err != nil {
				break
			}
			if err = w.Write(columnSlice); err != nil {
				log.Fatal(err)
			}
			counter++
			w.Flush()
		}
		if err = fs.Close(); err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()
	if err = nfs.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("success combine %d to %s\n", counter, csvPath)
}

func convertFeaturesByBorders(csvPath string) {
	ext := []string{"_ID3.csv", "_C4.5.csv", "_CART.csv"}
	for i := range ext {
		baseName := path.Base(csvPath)
		targetFile := filepath.Dir(csvPath) + "/" + baseName[:len(baseName)-4] + ext[i]
		nfs, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}
		w := csv.NewWriter(nfs)
		w.Comma = ','
		w.UseCRLF = true
		var counter int
		fs, err := os.Open(csvPath)
		if err != nil {
			log.Fatal(err)
		}
		r := csv.NewReader(fs)
		columnSlice, _ := r.Read()
		if err = w.Write(columnSlice[1:]); err != nil {
			log.Fatal(err)
		}
		counter++
		for {
			columnSlice, err = r.Read()
			if err != nil {
				break
			}
			columnSlice = columnSlice[1:]
			for j := range Config.Borders[i] {
				t, _ := strconv.ParseFloat(columnSlice[j], 64)
				if t > Config.Borders[i][j] {
					columnSlice[j] = "1"
				} else {
					columnSlice[j] = "0"
				}
			}
			if err = w.Write(columnSlice); err != nil {
				log.Fatal(err)
			}
			counter++
			w.Flush()
		}
		if err = fs.Close(); err != nil {
			log.Fatal(err)
		}

		w.Flush()
		if err = nfs.Close(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("success convert %d to %s\n", counter, targetFile)
	}
}

func getFlowDict(filePath string) map[string]struct{} {
	justID := make(map[string]struct{})
	fs, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(fs)
	r.Read()
	for {
		row, err := r.Read()
		if err != nil {
			break
		}
		justID[row[0]] = struct{}{}
	}
	return justID
}

func getAllColumnDict(filePath string) map[[CSVColumnNum]string]struct{} {
	hash := make(map[[CSVColumnNum]string]struct{})
	fs, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(fs)
	r.Read()
	var key [CSVColumnNum]string
	for {
		row, err := r.Read()
		if err != nil {
			break
		}
		for i := range key {
			key[i] = row[i]
		}
		hash[key] = struct{}{}
	}
	return hash
}
