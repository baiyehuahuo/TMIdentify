package calculate

import "strings"

func CountFreq(payloads []string, sigSet map[string]int) map[string]int {
	for sig := range sigSet {
		sigSet[sig] = OneSigCount(payloads, sig)
	}
	return sigSet
}

func OneSigCount(payloads []string, sig string) int {
	result := 0
	for i := range payloads {
		if strings.Contains(payloads[i], sig) {
			result++
		}
	}
	return result
}
