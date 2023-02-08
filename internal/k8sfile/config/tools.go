package config

import (
	"os"

	//"NBServer/tools/deploymentgen/template"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

type InnerConfig struct {
	OutputPath   string
	FilePath     string
	TemplatePath string
	Services     []map[interface{}]interface{}
	Default      ContentConfig
	Args         Args
}

// MustLoad
// loads yaml file into InnerConfig from path, exits on error.
func MustLoad(arg *Args) InnerConfig {

	c, err := LoadConfig(arg)
	if err != nil {
		log.Fatalf("error: config file %s, %s", arg.ConfigFile, err.Error())
	}
	c.TemplatePath = arg.Template
	c.OutputPath = arg.Dir
	c.FilePath = arg.Out
	return c
}

// LoadConfig
// loads yaml file into InnerConfig from path
func LoadConfig(arg *Args) (InnerConfig, error) {

	ic := InnerConfig{Args: *arg}
	var content []byte
	var err error

	if arg.JsonStr != "" {
		fmt.Println("Reading json Data")
		content = []byte(arg.JsonStr)
	} else {
		fmt.Println("Reading a Configuration File")
		content, err = os.ReadFile(arg.ConfigFile)
		//content, err = template.TempFs.ReadFile(arg.ConfigFile)

	}

	if err != nil {
		return ic, err
	}

	// 原始 config
	err = yaml.Unmarshal(content, &ic)
	if err != nil {
		fmt.Println("-->", err)
		return ic, err
	}

	// 将特定服务配置覆盖默认值
	for index, item := range ic.Services {

		defaultCopy := ic.Default

		bytes, err := yaml.Marshal(&item)
		if err != nil {
			return ic, err
		}

		err = yaml.Unmarshal(bytes, &defaultCopy)
		if err != nil {
			return ic, err
		}

		// 补充数据
		for portIndex, portItem := range defaultCopy.Ports {
			if portItem.TargetPort == 0 {
				defaultCopy.Ports[portIndex].TargetPort = portItem.Port
			}
		}

		// fixme: map、slice 会丢失数据

		// 生成 map[string]interface{}

		bytes, err = yaml.Marshal(&defaultCopy)
		if err != nil {
			return ic, err
		}
		ic.Services[index] = make(map[interface{}]interface{})
		err = yaml.Unmarshal(bytes, &ic.Services[index])
		if err != nil {
			return ic, err
		}
	}

	log.Println("load config success")

	return ic, nil
}

func UnmarshalToMap(content []byte) (map[interface{}]interface{}, error) {
	result := make(map[interface{}]interface{})

	// 注，原输入中有 []，不应该被转换
	str := string(content)
	str = strings.ReplaceAll(str, "\"", "#(0)")
	str = strings.ReplaceAll(str, "[", "'#(1)")
	str = strings.ReplaceAll(str, "]", "#(2)'")
	content = []byte(str)

	err := yaml.Unmarshal(content, &result)
	return result, err
}

func Marshal(content *map[interface{}]interface{}) ([]byte, error) {

	b, err := yaml.Marshal(content)
	if err != nil {
		return b, err
	}

	// 特殊处理
	str := string(b)
	str = strings.ReplaceAll(str, "#(2)'", "]")
	str = strings.ReplaceAll(str, "'#(1)", "[")
	str = strings.ReplaceAll(str, "#(0)", "\"")
	str = strings.ReplaceAll(str, "containerport", "containerPort")
	b = []byte(str)
	return b, nil
}
