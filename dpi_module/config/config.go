package config

import (
	"dpi_module/consts"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type GlobalConfigEntity struct {
	Mode               int     `yaml:"mode"`
	CsvPath            string  `yaml:"csv_path"`
	DirPath            string  `yaml:"dir_path"`
	OtherPath          string  `yaml:"other_path"`
	JudgeThreshold     int     `yaml:"judge_threshold"`
	SigMinLength       int     `yaml:"sig_min_length"`
	BaseThreshold      float64 `yaml:"base_threshold"`
	InclusionThreshold float64 `yaml:"inclusion_threshold"`
	PrintRawBytes      bool    `yaml:"print_raw_bytes"`
	PrintSigBytes      bool    `yaml:"print_sig_bytes"`
	SaveSigs           bool    `yaml:"save_sigs"`
	ReadSigPath        string  `yaml:"read_sig_path"`
	SaveSigName        string  `yaml:"save_sig_name"`
	MinPrecision       float64 `yaml:"min_precision"`
}

var GlobalConfig = &GlobalConfigEntity{}

func init() {
	yamlFile, err := ioutil.ReadFile(consts.DefaultYamlPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, GlobalConfig)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GlobalConfig read:  ", GlobalConfig)
}
