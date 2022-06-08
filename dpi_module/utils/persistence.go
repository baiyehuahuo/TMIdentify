package utils

import (
	"dpi_module/config"
	"dpi_module/consts"
	"dpi_module/types"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path"
)

func SaveSigs(sigMap map[string]float64, savePath string) {
	file, err := os.Create(savePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	enc := gob.NewEncoder(file)
	if err = enc.Encode(sigMap); err != nil {
		fmt.Println(err)
	}
}

func LoadSigs(loadPath string) map[string]float64 {
	sigMap := map[string]float64{}
	file, err := os.Open(loadPath)
	if err != nil {
		log.Fatal(err)
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&sigMap)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return sigMap
}

func SaveFlows(dirPath string, flowMap map[string]*types.Flow, flowPayloads []string) {
	mapFilePath := path.Join(dirPath, fmt.Sprintf("%s_%d", consts.DefaultFlowMap, config.GlobalConfig.JudgeThreshold))
	file, err := os.Create(mapFilePath)
	if err != nil {
		os.Remove(mapFilePath)
		fmt.Println(err)
		return
	}
	enc := gob.NewEncoder(file)
	if err = enc.Encode(flowMap); err != nil {
		os.Remove(mapFilePath)
		fmt.Println(err)
		return
	}

	payloadsFilePath := path.Join(dirPath, fmt.Sprintf("%s_%d", consts.DefaultFlowPayloads, config.GlobalConfig.JudgeThreshold))
	file, err = os.Create(payloadsFilePath)
	if err != nil {
		os.Remove(mapFilePath)
		os.Remove(payloadsFilePath)
		fmt.Println(err)
		return
	}
	enc = gob.NewEncoder(file)
	if err = enc.Encode(flowPayloads); err != nil {
		os.Remove(mapFilePath)
		os.Remove(payloadsFilePath)
		fmt.Println(err)
		return
	}
}

func LoadFlows(dirPath string) (flowMap map[string]*types.Flow, flowPayloads []string) {
	flowMap = map[string]*types.Flow{}
	file, err := os.Open(path.Join(dirPath, fmt.Sprintf("%s_%d", consts.DefaultFlowMap, config.GlobalConfig.JudgeThreshold)))
	if err != nil {
		return nil, nil
	}
	dec := gob.NewDecoder(file)
	if err = dec.Decode(&flowMap); err != nil {
		return nil, nil
	}
	file, err = os.Open(path.Join(dirPath, fmt.Sprintf("%s_%d", consts.DefaultFlowPayloads, config.GlobalConfig.JudgeThreshold)))
	if err != nil {
		return nil, nil
	}
	dec = gob.NewDecoder(file)
	if err = dec.Decode(&flowPayloads); err != nil {
		return nil, nil
	}
	return flowMap, flowPayloads
}

func GetRawSaveName() string {
	gConfig := config.GlobalConfig
	return path.Join(gConfig.DirPath, fmt.Sprintf("%s_%d_%.2f_%.2f", gConfig.SaveSigName, gConfig.JudgeThreshold, gConfig.BaseThreshold, gConfig.InclusionThreshold))
}

func GetFilterSaveName() string {
	gConfig := config.GlobalConfig
	return path.Join(gConfig.DirPath, fmt.Sprintf("%s_%d_%.2f_%.2f_%.2f", gConfig.SaveSigName, gConfig.JudgeThreshold, gConfig.BaseThreshold, gConfig.InclusionThreshold, gConfig.MinPrecision))
}
