package config

import (
	"discretization/consts"
	"io/ioutil"
	"log"
	"path"

	"gopkg.in/yaml.v2"
)

type GlobalConfigEntity struct {
	Mode                  int         `yaml:"mode"`
	MaxDeeper             int         `yaml:"max_deeper"`
	CSVFilePath           string      `yaml:"csv_file_path"`
	Borders               [][]float64 `yaml:"borders"`
	UseInputBorders       bool        `yaml:"use_input_borders"`
	GetBetterRebuildTimes int         `yaml:"get_better_rebuild_times"`
	TestRate              int         `yaml:"test_rate"`
	TestModelTimes        int         `yaml:"test_model_times"`
	SaveModel             bool        `yaml:"save_model"`
	SavePath              []string    `yaml:"model_save_path"`
	ModelPath             []string    `yaml:"model_path"`
	PrintBestGain         bool        `yaml:"print_best_gain"`
	PrintCalBorder        bool        `yaml:"print_cal_border"`
	PrintMatrix           bool        `yaml:"print_matrix"`
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

	// if len(GlobalConfig.SavePath) == 0 {
	// 	GlobalConfig.SavePath = CreateModelPath(GlobalConfig.CSVFilePath)
	// }
	// if len(GlobalConfig.ModelPath) == 0 {
	// 	GlobalConfig.ModelPath = CreateModelPath(GlobalConfig.CSVFilePath)
	// }

	// fmt.Println("GlobalConfig read.", GlobalConfig)
}

func CreateModelPath(csvFilePath string) (modelPaths []string) {
	csvFileName := path.Base(csvFilePath)
	modelBaseName := csvFileName[:len(csvFileName)-len(path.Ext(csvFileName))]
	for _, suffix := range consts.ModeSuffix {
		modelPaths = append(modelPaths, consts.ModelDir+"/"+modelBaseName+suffix)
	}
	return modelPaths
}
