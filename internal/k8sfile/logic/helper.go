package logic

import (
	"NBServer/tools/deploymentgen/internal/k8sfile/config"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	DefaultTemplatePath = "template/"
	DefaultOutputPath   = ""
	OutputFileName      = "resources"
	OutputSuffix        = "-kube.yaml"
	TemplateSuffix      = ".yaml"
	InnerSplit          = "\n---\n\n"
	KeySplit            = "."

	DeploymentName     = "deployment"
	ServiceName        = "service"
	ServiceNetworkName = "service-network"
	IngressName        = "ingress"
	HPAName            = "hpa"

	MarkTemplateStart = "$("
	MarkTemplateEnd   = ")"
)

type templateConfig struct {
	templatePath string
	outputPath   string
	fileName     string
	content      map[interface{}]interface{}
}

var (
	globalIndex = 0
)

// 输出yaml文件
func GenerateServicesYaml(c *config.InnerConfig) {

	tc := templateConfig{
		templatePath: getRealValue(DefaultTemplatePath, getPath(c.TemplatePath)),
		outputPath:   getRealValue(DefaultOutputPath, getPath(c.OutputPath)),
		fileName:     getRealValue(OutputFileName, c.FilePath),
	}

	if tc.outputPath != "" {
		mkdir(tc.outputPath)
	}

	for _, item := range c.Services {

		// 修改 image 数据
		if c.Args.Tag != "" {
			var image = replaceImageTag(item["image"].(string), c.Args.Tag)
			item["image"] = image
		}

		generateServiceYaml(&tc, item)
	}
}

func generateServiceYaml(tc *templateConfig, c map[interface{}]interface{}) {
	if tc == nil || c == nil {
		log.Fatalf("Unknown config input")
	}

	//name, ok := c["name"]
	//if !ok {
	//	log.Fatalf("undefined service name")
	//}

	fileName := ""
	if tc.outputPath == "" {
		fileName = fmt.Sprintf("%s%s", tc.fileName, OutputSuffix)
	} else {
		fileName = fmt.Sprintf("%s%s%s", tc.outputPath, tc.fileName, OutputSuffix)
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		log.Fatalf("cannot create output file, error: %s", err.Error())
	}
	defer file.Close()

	// 将 map 的深度打散，格式为 key1.key2;
	// 数组的key为0、1、2...
	// 保留关键字 key, key 会单独索引
	keyMap := make(map[string]interface{})

	regenerateMap(keyMap, c, "")

	tc.content = getContent(tc.templatePath, DeploymentName)

	// 将配置文件的信息替换到模板中
	handleContent(tc.content, keyMap, file, false)

	if _, ok := keyMap["ports.network"]; ok {
		// 使用 network service
		tc.content = getContent(tc.templatePath, ServiceNetworkName)
		handleContent(tc.content, keyMap, file, true)
	} else {
		// 使用默认 service
		tc.content = getContent(tc.templatePath, ServiceName)
		handleContent(tc.content, keyMap, file, true)
	}

	if keyMap["host"] != "" {
		//添加了host字段，就会生成ingress yaml
		// 使用 ingress
		tc.content = getContent(tc.templatePath, IngressName)
		handleContent(tc.content, keyMap, file, true)
	}

	if value, ok := keyMap["hpa"]; ok {
		valueSlice := value.([]interface{})

		for hpaIndex := range valueSlice {
			globalIndex = hpaIndex
			tc.content = getContent(tc.templatePath, HPAName)
			handleContent(tc.content, keyMap, file, true)
		}
	}
}

// 模板   数据  文件
func handleContent(content map[interface{}]interface{}, c map[string]interface{}, file *os.File, useSplit bool) {
	if useSplit {
		_, err := file.Write([]byte(InnerSplit))
		if err != nil {
			log.Fatalf("cannot write output file, error: %s", err.Error())
		}
	}

	handleContentMap(content, c)

	bytes, err := config.Marshal(&content)

	if err != nil {
		log.Fatalf("Unknown yaml format")
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Fatalf("cannot write output file, error: %s", err.Error())
	}

}

// 打散 map
func regenerateMap(result map[string]interface{}, item interface{}, k string) {
	// 判断数据类型
	vTypeKind := reflect.TypeOf(item).Kind()

	switch vTypeKind {
	case reflect.Map:
		// 递归
		valueMap := item.(map[interface{}]interface{})
		for key, value := range valueMap {
			var newK string
			if k == "" {
				newK = key.(string)
			} else {
				newK = k + KeySplit + key.(string)
			}
			result[newK] = value

			regenerateMap(result, value, newK)
		}
		break
	case reflect.Slice:
		// 递归
		valueSlice := item.([]interface{})
		for key, value := range valueSlice {
			var newK string
			if k == "" {
				newK = strconv.Itoa(key)
			} else {
				newK = k + KeySplit + strconv.Itoa(key)
			}
			result[newK] = value
			regenerateMap(result, value, newK)

			// type 特殊处理
			if reflect.TypeOf(value).Kind() == reflect.Map {
				valueMap := value.(map[interface{}]interface{})
				if valueItem, ok := valueMap["type"]; ok {
					s := valueItem.(string)
					if s != "" {
						newK = k + KeySplit + s
						result[newK] = value
						regenerateMap(result, value, newK)
					}
				}
			}
		}
		break
	}
}

// 替换 map
func handleContentMap(content map[interface{}]interface{}, c map[string]interface{}) {
	for key, value := range content {

		handleValue(content, key, value, c)
	}
}

// 替换 slice
func handleContentSlice(content []interface{}, c map[string]interface{}) {
	for key, value := range content {
		handleValue(content, key, value, c)
	}
}

func handleValue(content interface{}, key interface{}, value interface{}, c map[string]interface{}) {
	if value == nil {
		return
	}

	vTypeKind := reflect.TypeOf(value).Kind()

	switch vTypeKind {
	case reflect.String:
		valueStr := value.(string)
		realValue := replaceMark(valueStr, c)

		cTypeKind := reflect.TypeOf(content).Kind()
		if cTypeKind == reflect.Map {
			cMap := content.(map[interface{}]interface{})
			cMap[key] = realValue
		} else if cTypeKind == reflect.Slice {
			cSlice := content.([]interface{})
			cSlice[key.(int)] = realValue
		}
		break
	case reflect.Map:
		// 递归
		valueMap := value.(map[interface{}]interface{})
		handleContentMap(valueMap, c)
		break
	case reflect.Slice:
		// 递归
		valueSlice := value.([]interface{})
		handleContentSlice(valueSlice, c)
		break
	}
}

type PortConfig struct {
	ContainerPort int
}

func replaceMark(str string, c map[string]interface{}) interface{} {
	for {
		indexStart := strings.Index(str, MarkTemplateStart)
		indexEnd := strings.Index(str, MarkTemplateEnd)
		if indexStart < 0 || indexEnd < 0 {
			return str
		}

		dictKey := strings.ReplaceAll(str[indexStart+2:indexEnd], "#", strconv.Itoa(globalIndex))

		if value, ok := c[dictKey]; ok {
			vType := reflect.TypeOf(value).Kind()
			switch vType {
			case reflect.String:
				str = strings.Replace(str, str[indexStart:indexEnd+1], value.(string), 1)
				break
			default:
				// HACK: 需要使用 handler 来处理
				if dictKey == "ports" {
					realSlice := value.([]interface{})
					result := make([]PortConfig, len(realSlice))
					for index, item := range realSlice {
						containerPort := item.(map[interface{}]interface{})["port"].(int)

						result[index] = PortConfig{
							containerPort,
						}
					}
					return result
				}

				// 其他值无法拼接，直接返回原值
				return value
			}
		} else {
			// 不存在，直接去除无效标记
			str = strings.Replace(str, str[indexStart:indexEnd+1], "", 1)
		}
	}
}
