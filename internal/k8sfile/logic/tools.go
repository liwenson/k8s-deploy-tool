package logic

import (
	"NBServer/tools/deploymentgen/internal/k8sfile/config"
	"NBServer/tools/deploymentgen/template"
	"fmt"
	"log"
	"os"
	"strings"
)

func getPath(value string) string {
	if value == "" || value[len(value)-1:] == "/" {
		return value
	}
	return value + "/"
}

func getRealValue(defaultValue, inputValue string) string {

	if inputValue == "" {
		return defaultValue
	}
	return inputValue
}

func getContent(path, name string) map[interface{}]interface{} {

	exists, err := pathExists(path)
	if err != nil {
		log.Fatalf("template file %s not find", err)
	}

	var fullPath string
	var content []byte
	if exists {
		fullPath = fmt.Sprintf("%s%s%s", path, name, TemplateSuffix)
		content, err = os.ReadFile(fullPath)
		if err != nil {
			log.Fatalf("template %s not find", name)
		}
	} else {
		fullPath = fmt.Sprintf("%s%s", name, TemplateSuffix)
		content, err = template.TempFs.ReadFile(fullPath)
		if err != nil {
			log.Fatalf("template %s not find", name)
		}
	}

	result, err := config.UnmarshalToMap(content)
	if err != nil {
		log.Fatalf("template %s invalid, err: %s", name, err.Error())
	}

	return result
}

// mkdir 没有则创建文件夹，有则清除所有文件
func mkdir(path string) {
	_, err := os.Stat(path)
	if err == nil {
		err = os.RemoveAll(path)
		if err != nil {
			log.Fatalf("cannot clear old file, error: %s", err.Error())
		}
	}
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		log.Fatalf("output folder create failed")
	}
}

func replaceImageTag(image, tag string) string {
	countSplit := strings.Split(image, ":")
	if len(countSplit) == 2 {
		countSplit[1] = tag
		return fmt.Sprintf("%s:%s", countSplit[0], tag)
	}
	return fmt.Sprintf("%s:%s", image, tag)

}

// pathExists 判断所给路径文件/文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}
