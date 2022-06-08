package config

import (
	"encoding/csv"
	"fmt"
	"gopkg.in/yaml.v2"
	"identifyTencentMeeting/consts"
	"io/ioutil"
	"log"
	"os"
)

type GlobalConfigEntity struct {
	TMPath             string   `yaml:"tm_path"`
	OtherPath          string   `yaml:"other_path"`
	UseDPI             bool     `yaml:"use_dpi"`
	UseML              bool     `yaml:"use_ml"`
	MLJudgePacket      int      `yaml:"ml_judge_packet"`
	DPIJudgePacket     int      `yaml:"dpi_judge_packet"`
	DPIBaseThreshold   float64  `yaml:"dpi_base_threshold"`
	DPIFilterPrecision float64  `yaml:"dpi_filter_precision"`
	SplitFlow          bool     `yaml:"split"`
	FilterShort        bool     `yaml:"filter_short"`
	SigPath            string   `yaml:"sig_path"`
	TreesPath          []string `yaml:"trees_path"`
	PrintDetect        bool     `yaml:"print_detect"`
	Label              bool     `yaml:"label"`
	HostIP             []string
	CSVFile            *os.File
	CSVWriter          *csv.Writer
	CSVColumn          int
}

var GlobalConfig = &GlobalConfigEntity{}

func (g GlobalConfigEntity) isSrcIP(ip string) bool {
	for i := range g.HostIP {
		if g.HostIP[i] == ip {
			return true
		}
	}
	return false
}

func init() {
	yamlFile, err := ioutil.ReadFile(consts.DefaultYamlPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, GlobalConfig)
	if err != nil {
		log.Fatal(err)
	}

	UpdatePath()
	fmt.Println("GlobalConfig read.", GlobalConfig)
}

func UpdatePath() {
	sigPath := "classify/DPI/TM_%d_%.2f_0.50_%.2f"
	treePath := []string{
		"classify/ML/decision_tree/NonVPN_all_%d_with_TencentMeeting_ID3",
		"classify/ML/decision_tree/NonVPN_all_%d_with_TencentMeeting_C4.5",
		"classify/ML/decision_tree/NonVPN_all_%d_with_TencentMeeting_CART",
	}
	GlobalConfig.SigPath = fmt.Sprintf(sigPath, GlobalConfig.DPIJudgePacket, GlobalConfig.DPIBaseThreshold, GlobalConfig.DPIFilterPrecision)
	GlobalConfig.TreesPath = nil
	for i := range treePath {
		GlobalConfig.TreesPath = append(GlobalConfig.TreesPath, fmt.Sprintf(treePath[i], GlobalConfig.MLJudgePacket))
	}
}
