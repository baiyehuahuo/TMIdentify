package DPI

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"
)

var sigMap map[string]float64

func LoadSigMap(sigPath string) error {
	sigMap = map[string]float64{}
	file, err := os.Open(sigPath)
	if err != nil {
		log.Fatal(err)
	}
	dec := gob.NewDecoder(file)
	return dec.Decode(&sigMap)
}

func Classify(flowPayload string) bool {
	// fmt.Println("DPI Classify.")
	for sig := range sigMap {
		if strings.Contains(flowPayload, sig) {
			// fmt.Printf("sig: %s is matched, %v.\n", sig, []byte(sig))
			return true
		}
	}
	return false
}

func PrintSigs() {
	for sig := range sigMap {
		fmt.Println([]byte(sig))
	}
}
