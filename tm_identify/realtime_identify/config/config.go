package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"identifyTencentMeeting/consts"
	"io/ioutil"
	"log"
)

type GlobalConfigEntity struct {
	DeviceName     string   `yaml:"device_name"`
	MLJudgePacket  int      `yaml:"ml_judge_packet"`
	DPIJudgePacket int      `yaml:"dpi_judge_packet"`
	SigPath        string   `yaml:"sig_path"`
	TreesPath      []string `yaml:"trees_path"`
	HostIP         []string
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

	fmt.Println("GlobalConfig read.")
}
