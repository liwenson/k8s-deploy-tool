package logic

import (
	"NBServer/tools/deploymentgen/internal/getgitlabfile/config"
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

var Key = []byte("abcdabcdabcdabcdabcdabcdabcdabcd")

func generalWrite(filename string, param []byte) {
	/*
		O_RDONLY	打开只读文件
		O_WRONLY	打开只写文件
		O_RDWR	打开既可以读取又可以写入文件
		O_APPEND	写入文件时将数据追加到文件尾部
		O_CREATE	如果文件不存在，则创建一个新的文件
		O_EXCL	文件必须不存在，然后会创建一个新的文件
		O_SYNC	打开同步I/0
		O_TRUNC	文件打开时可以截断

	*/

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_SYNC|os.O_TRUNC|os.O_RDWR, 0766)
	if err != nil {
		log.Fatal("open file error :", err)
		return
	}

	// 关闭文件
	defer f.Close()

	write := bufio.NewWriter(f)

	// 关闭文件
	defer f.Close()

	// 字节方式写入
	_, err = write.Write(param)
	if err != nil {
		log.Fatal(err)
		return
	}

	//Flush将缓存的文件真正写入到文件中
	err = write.Flush()
	if err != nil {
		log.Fatal(err)
		return
	}

	time.Sleep(time.Duration(1) * time.Second)

	// 字节方式写入
	//_, err = f.Write(param)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
}

// 调用os.MkdirAll递归创建文件夹
func createDir(filePath string) error {
	if !isExist(filePath) {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// mkdir 没有则创建文件夹,
func mkdir(path string) {
	_, err := os.Stat(path)
	if err == nil {
		return
	}
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		log.Fatalf("output folder create failed")
	}
}

// mkdir 没有则创建目录,有则清空目录
func mkdirClean(path string) {
	_, err := os.Stat(path)
	if err == nil {
		return
	}
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		log.Fatalf("output folder create failed")
	}
}

func PwdEncryption(srcData string) {

	//测试加密
	encData, err := ECBEncrypt([]byte(srcData), Key)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(encData))

}

func PwdDecrypt(encData string) string {

	encData = strings.Trim(encData, " ")

	arr, _ := base64.StdEncoding.DecodeString(encData)

	//测试解密
	decData, err := ECBDecrypt(arr, Key)
	if err != nil {
		panic(err.Error())
	}

	return string(decData)
}

func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func FileTree(git *gitlab.Client, data []*gitlab.TreeNode, args config.Args, result []*gitlab.TreeNode) []*gitlab.TreeNode {

	for _, re := range data {

		if re.Type == "blob" {
			result = append(result, re)

		} else {
			branch := new(string)
			*branch = args.Branch

			sp := new(string)
			*sp = re.Path

			listRess := gitlab.ListTreeOptions{
				Path: sp,
				Ref:  branch,
			}
			ress, _, err := git.Repositories.ListTree(args.Pid, &listRess)
			if err != nil {
				log.Fatal(err)
			}

			res := FileTree(git, ress, args, result)
			for _, node := range res {
				result = append(result, node)
			}
		}
	}
	return result
}

func createPath(base string, arrays []string, args config.Args) string {
	// 删除第一个元素
	// 删除最后一个元素
	// 拼接路径

	arr := []string{}

	arr = arrays[:len(arrays)-1]

	var build strings.Builder
	build.WriteString(base)

	for i, _ := range arr {

		if !config.Keepdirectory && args.Path == arr[i] {
			continue
		}
		build.WriteString("/")
		build.WriteString(arr[i])
	}

	dirname := build.String()

	if args.Out == "" {
		dirname = strings.TrimLeft(dirname, "/")
	}

	err := createDir(dirname)
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	return dirname
}
