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
		fmt.Println(err)
		log.Fatal(err)
	}
	dec := gob.NewDecoder(file)

	return dec.Decode(&sigMap)
}

func Classify(flowPayload string) bool {
	for sig := range sigMap {
		if strings.Contains(flowPayload, sig) {
			// fmt.Printf("sig: %s is matched, %v.\n", sig, []byte(sig))
			return true
		}
	}
	return false
}
