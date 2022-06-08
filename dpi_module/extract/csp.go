package extract

import (
	"dpi_module/calculate"
	"dpi_module/config"
	"fmt"
	"strings"
)

func CSP(payloads []string, minLength int, baseThreshold, inclusionThreshold float64) []string {
	sigSets := make([]map[string]int, 0, minLength+1)
	thresholdFreq := int(float64(len(payloads)) * baseThreshold)
	sigSets = append(sigSets, map[string]int{"": len(payloads)})
	var sigSet = map[string]int{}
	for i := range payloads {
		for j := range payloads[i] {
			sigSet[payloads[i][j:j+1]]++
		}
	}
	sigSet = filterByFreq(sigSet, thresholdFreq)
	sigSets = append(sigSets, sigSet)
	sigLength := 1
	for len(sigSets[sigLength]) > 0 {
		sigSet = getNextRawSet(sigSet, sigLength)
		sigSet = calculate.CountFreq(payloads, sigSet)
		sigSet = filterByFreq(sigSet, thresholdFreq)
		fmt.Printf("sigLength: %v, len(sigSet): %d\n", sigLength+1, len(sigSet))
		sigSets = append(sigSets, sigSet)
		sigLength++
	}
	maxSigLength := sigLength - 1
	if maxSigLength < minLength {
		return nil
	}

	result := make([]string, 0, len(sigSets[maxSigLength]))
	var i, j, length int
	for i = maxSigLength; i >= minLength; i-- {
		length = len(result)
		thresholdFreq = int(float64(len(payloads)) * (inclusionThreshold - 0.01*float64(i-config.GlobalConfig.SigMinLength)))
		for sig, freq := range sigSets[i] {
			if freq < thresholdFreq {
				continue
			}
			for j = 0; j < length; j++ {
				if strings.Contains(result[j], sig) {
					break
				}
			}
			if j == length {
				result = append(result, sig)
			}
		}
	}

	resultCount(payloads, result)
	return result
}

func resultCount(payloads []string, sigs []string) {
	sigSet := make(map[string]int, len(sigs))
	for _, sig := range sigs {
		sigSet[sig] = 0
	}
	sigSet = calculate.CountFreq(payloads, sigSet)
	for sig, val := range sigSet {
		fmt.Printf("freqPercent: %v\tfreq: %v\t sigLength:%v \n", float64(val)/float64(len(payloads)), val, len(sig))
		byteSigs := fmt.Sprintf("%v", []byte(sig))
		byteSigs = strings.ReplaceAll(byteSigs, " ", ", ")
		byteSigs = "{" + byteSigs[1:len(byteSigs)-1] + "}"
		fmt.Println(byteSigs, sig)
	}
}

func getNextRawSet(sigSet map[string]int, length int) map[string]int {
	nextSigSet := make(map[string]int)
	for key1 := range sigSet {
		for key2 := range sigSet {
			// fmt.Println(key1, key2, key1[1:], key2[:length-1], key1[1:] == key2[:length-1])
			if key1[1:] == key2[:length-1] {
				// fmt.Println(length, key1[0:1]+key2)
				nextSigSet[key1[0:1]+key2] = 0
			}
		}
	}
	return nextSigSet
}

func filterByFreq(sigSet map[string]int, thresholdFreq int) map[string]int {
	for sig, freq := range sigSet {
		if freq < thresholdFreq {
			delete(sigSet, sig)
		}
	}
	return sigSet
}
