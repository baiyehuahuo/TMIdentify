package main

import (
	"dpi_module/calculate"
	"dpi_module/config"
	"dpi_module/extract"
	"dpi_module/types"
	"dpi_module/utils"
	"fmt"
	"strings"
)

func cspExtract() {
	gConfig := config.GlobalConfig
	flowCount := types.GetTotal()
	payloads := types.GetFlowPayloads()
	if gConfig.PrintRawBytes {
		utils.PrintStringsByByte(payloads)
	}
	fmt.Printf("TotalSatisfyLenPayloads: %d \t total flow: %d\n", len(payloads), flowCount)
	cspSigs := extract.CSP(payloads, gConfig.SigMinLength, gConfig.BaseThreshold, gConfig.InclusionThreshold)
	fmt.Printf("total cspSigs: %d\ttotal Payloads: %d\n", len(cspSigs), len(payloads))

	if gConfig.SaveSigs {
		sigMap := make(map[string]float64, len(cspSigs))
		for _, sig := range cspSigs {
			sigMap[sig] = 0
		}
		utils.SaveSigs(sigMap, utils.GetRawSaveName())
	}
}

func calculatePrecision(sigMap map[string]float64, targetPcapDir string, otherPcapDir string) {
	gConfig := config.GlobalConfig
	TPMap, FPMap := make(map[string]int, len(sigMap)), make(map[string]int, len(sigMap))
	for sig := range sigMap {
		TPMap[sig] = 0
		FPMap[sig] = 0
	}

	utils.ReadPcapDir(targetPcapDir)
	payloads := types.GetFlowPayloads()
	types.Clean()
	TPMap = calculate.CountFreq(payloads, TPMap)

	utils.ReadPcapDir(otherPcapDir)
	payloads = types.GetFlowPayloads()
	types.Clean()
	FPMap = calculate.CountFreq(payloads, FPMap)

	for sig := range sigMap {
		sigMap[sig] = float64(TPMap[sig]) / float64(TPMap[sig]+FPMap[sig])
		if sigMap[sig] < gConfig.MinPrecision || allZeroOne(sig) {
			delete(sigMap, sig)
		} else {
			// fmt.Println(TPMap[sig], FPMap[sig])
		}
	}
	if gConfig.SaveSigs {
		utils.SaveSigs(sigMap, utils.GetFilterSaveName())
	}
}

func allZeroOne(sig string) bool {
	for i := range sig {
		if sig[i] != 0 && sig[i] != 1 {
			return false
		}
	}
	return true
}

func totalEvaluate(sigMap map[string]float64, targetPcapDir string, otherPcapDir string) {
	TP, FP, TN, FN := 0, 0, 0, 0
	totalCountMap, unrecognizedCountMap := make(map[int]int), make(map[int]int)
	utils.ReadPcapDir(targetPcapDir)
	payloads := types.GetFlowPayloads()
	types.Clean()
	for _, payload := range payloads {
		totalCountMap[len(payload)]++
		if classifyByMap(sigMap, payload) {
			TP++
		} else {
			// fmt.Println([]byte(payload))
			unrecognizedCountMap[len(payload)]++
			FN++
		}
	}
	for key := range totalCountMap {
		if _, ok := unrecognizedCountMap[key]; !ok {
			unrecognizedCountMap[key] = 0
		}
	}
	utils.PrintSplitCountMap(totalCountMap, "total")
	utils.PrintSplitCountMap(unrecognizedCountMap, "unrecognized")
	fmt.Println("total tencent meeting payloads: ", len(payloads))
	utils.ReadPcapDir(otherPcapDir)
	payloads = types.GetFlowPayloads()
	types.Clean()
	for _, payload := range payloads {
		if classifyByMap(sigMap, payload) {
			for sig := range sigMap {
				if strings.Contains(payload, sig) {
					fmt.Println([]byte(sig))
				}
			}
			FP++
		} else {
			TN++
		}
	}
	fmt.Println("total other payloads: ", len(payloads))
	fmt.Printf("TP: %d\tTN: %d\tFP: %d\t FN: %d\n", TP, TN, FP, FN)
	accuracy := float64(TP+TN) / float64(TP+TN+FP+FN)
	precision := float64(TP) / float64(TP+FP)
	recall := float64(TP) / float64(TP+FN)
	f1 := (2 * precision * recall) / (precision + recall)
	fmt.Printf("accuracy: %f\tprecision: %f\trecall: %f\tf1: %f\n", accuracy, precision, recall, f1)
}

func classifyByMap(sigMap map[string]float64, payload string) bool {
	for sig := range sigMap {
		if strings.Contains(payload, sig) {
			return true
		}
	}
	return false
}

func deleteSig(sigMap map[string]float64, sig string) {
	delete(sigMap, sig)
	utils.SaveSigs(sigMap, utils.GetFilterSaveName())
}

func addSig(sigMap map[string]float64, sig string) {
	sigMap[sig] = 1
	utils.SaveSigs(sigMap, utils.GetFilterSaveName())
}
