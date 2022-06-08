package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ConfigEntity struct {
	Mode            int         `yaml:"mode"`
	DicPath         string      `yaml:"dic_path"`
	TotalPath       string      `yaml:"total_path"`
	CSVPath         string      `yaml:"csv_path"`
	Label           bool        `yaml:"label"`
	DirPath         string      `yaml:"dir_path"`
	CombineCSVPath  string      `yaml:"combine_csv_path"`
	CombineCSVPaths []string    `yaml:"combine_csv_paths"`
	Borders         [][]float64 `yaml:"borders"`
}

var Config = &ConfigEntity{}

func init() {
	yamlFile, err := ioutil.ReadFile(YamlPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, Config)

	if err != nil {
		log.Fatal(err)
	}

	if Config.CombineCSVPath == "" {
		Config.CombineCSVPath = Config.DirPath + ".csv"
	}
	if Config.CSVPath == "" {
		Config.CSVPath = Config.DirPath + ".csv"
	}

	fmt.Println("GlobalConfig read.", Config)
}
