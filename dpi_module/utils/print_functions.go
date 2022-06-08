package utils

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func PrintStringsByByte(strs []string) {
	for _, str := range strs {
		fmt.Println([]byte(str))
	}
}

func PrintSigBytes(sigs map[string]float64) {
	// t := string([]byte{0, 0, 0, 13})
	length := 0
	for sig := range sigs {
		fmt.Println(sig, []byte(sig))
		length += len(sig)
		// if sig == t {
		// 	fmt.Println("???")
		// }
	}
	fmt.Println(float64(length)/float64(len(sigs)), len(sigs))
}

func PrintSplitCountMap(countMap map[int]int, label string) {
	fmt.Printf("%s: \n", label)
	type pair struct {
		Length, Count int
	}
	var pairs []pair
	for key, val := range countMap {
		pairs = append(pairs, pair{key, val})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Length < pairs[j].Length
	})
	lengths := make([]int, 0, len(pairs))
	counts := make([]int, 0, len(pairs))
	for i := range pairs {
		lengths = append(lengths, pairs[i].Length)
		counts = append(counts, pairs[i].Count)
	}
	fmt.Println(strings.ReplaceAll(fmt.Sprintf("%v", lengths), " ", ","))
	fmt.Println(strings.ReplaceAll(fmt.Sprintf("%v", counts), " ", ","))
}

func PrintSliceSigToString(str string) string {
	// fmt.Println(str)
	str = str[1 : len(str)-1]
	str = strings.ReplaceAll(str, "\n", "")
	sigByteString := strings.Split(str, " ")
	// fmt.Println(sigByteString)
	buffer := bytes.Buffer{}
	for i := range sigByteString {
		b, _ := strconv.Atoi(sigByteString[i])
		buffer.WriteByte(byte(b))
	}
	// fmt.Println([]byte(buffer.String()))
	return buffer.String()
}
