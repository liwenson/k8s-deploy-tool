package logic

import (
	"NBServer/tools/deploymentgen/internal/dockerfile/config"
	tpl "NBServer/tools/deploymentgen/template"
	"bufio"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"os"
)

func Gendockerfile(template config.DockerfileTemplate) {

	//context := pongo2.Context{
	//	"FromImg": "reg.test.com/library/openjdk-alpine-with-chinese-timezone:8-jdk",
	//	"AppName": "be-bbs-omp-asc",
	//	"Outputs": "./target/logdemo.jar",
	//	"JvmOpt":  "",
	//}

	var (
		filePath string
		tempTpl  *pongo2.Template
	)

	exists, err := pathExists(template.TplFile)
	if err != nil {
		fmt.Println("检查文件或目录是否存在发生错误", err)
	}

	if template.OutFile == "" {
		filePath = "Dockerfile"
	} else {
		filePath = template.OutFile
	}

	file, err := os.OpenFile(fmt.Sprintf("%s%s", "", filePath), os.O_WRONLY|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}

	//及时关闭file句柄
	defer file.Close()

	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)

	if exists {
		fmt.Println("Using template files：", template.TplFile)
		tempTpl = pongo2.Must(pongo2.FromFile(template.TplFile))
	} else {
		set := pongo2.NewSet("embed", pongo2.NewFSLoader(tpl.TempFs))
		tempTpl = pongo2.Must(set.FromFile("dockerfile.template"))
	}

	err = tempTpl.ExecuteWriter(template.Context, write)
	if err != nil {
		fmt.Print(err)
	}

	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	fmt.Println("Dockerfile generated successfully")
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
