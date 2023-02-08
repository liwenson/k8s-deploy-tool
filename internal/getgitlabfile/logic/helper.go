package logic

import (
	"NBServer/tools/deploymentgen/internal/getgitlabfile/config"
	"github.com/xanzy/go-gitlab"
	"log"
	"regexp"
	"strings"
)

func GitlabRun(args config.Args) {
	/**
	http://gitlab.ztoyc.com/api/v4
	FszEcpyWdCZb3b3XDDcw
	*/

	url := ""
	url = strings.TrimRight(args.Server, "/")

	re := regexp.MustCompile("(http|https)://")

	if !re.MatchString(url) {

		var build strings.Builder
		build.WriteString("http://")
		build.WriteString(url)
		build.WriteString("/api/v4")

		url = build.String()
	}

	git, err := gitlab.NewClient(args.Token, gitlab.WithBaseURL(url))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	gf := &gitlab.GetFileMetaDataOptions{
		Ref: gitlab.String(args.Branch),
	}

	// 判断路径是不是文件
	f, _, err := git.RepositoryFiles.GetFileMetaData(args.Pid, args.Path, gf)
	if f == nil && err != nil {

		// 处理目录文件下载
		res := GetFileTree(git, args)
		FileSave(git, res, args)

	} else {
		// 处理单文件下载
		res := GetFils(args)
		FileSave(git, res, args)

	}
	return

	//countSplit := strings.Split(args.Path, ".")
	//if len(countSplit) > 1 {
	//	// 处理单文件下载
	//	res := GetFils(args)
	//	FileSave(git, res, args)
	//	return
	//}
	//
	//res := GetFileTree(git, args)
	////log.Println("res ", res)
	//FileSave(git, res, args)

}

func GetFils(args config.Args) []*gitlab.TreeNode {

	countSplit := strings.Split(args.Path, "/")

	filename := countSplit[len(countSplit)-1]

	var node []*gitlab.TreeNode
	res := gitlab.TreeNode{
		Name: filename,
		Path: args.Path,
	}
	node = append(node, &res)

	return node

}

// 获取目录中所有的下载地址
func GetFilelist(git *gitlab.Client, args config.Args) []*gitlab.TreeNode {

	branch := new(string)
	*branch = args.Branch

	sp := new(string)
	*sp = args.Path

	listRess := gitlab.ListTreeOptions{
		Path: sp,
		Ref:  branch,
	}

	res, _, err := git.Repositories.ListTree(args.Pid, &listRess)

	if err != nil {
		log.Println(err)
	}
	return res
}

// 遍历目录中所有的下载地址
func GetFileTree(git *gitlab.Client, args config.Args) []*gitlab.TreeNode {
	branch := new(string)
	*branch = args.Branch

	sp := new(string)
	*sp = args.Path

	listRess := gitlab.ListTreeOptions{
		Path: sp,
		Ref:  branch,
	}

	res, _, err := git.Repositories.ListTree(args.Pid, &listRess)
	if err != nil {
		log.Println(err)
	}

	var (
		tns []*gitlab.TreeNode
	)

	press := FileTree(git, res, args, tns)

	for _, node := range press {
		log.Println("node: ", node.Path)
	}
	return press
}

// 文件保存
func FileSave(git *gitlab.Client, fis []*gitlab.TreeNode, args config.Args) {

	branch := new(string)
	*branch = args.Branch

	if len(fis) < 2 {
		file, _, err := git.RepositoryFiles.GetRawFile(args.Pid, fis[0].Path, &gitlab.GetRawFileOptions{
			Ref: branch,
		})

		if err != nil {
			log.Println(err)
			return
		}

		generalWrite(fis[0].Name, file)
		return
	}
	for _, fi := range fis {

		//_, _, err := git.RepositoryFiles.GetRawFile(args.Pid, fi.Path, &gitlab.GetRawFileOptions{
		file, _, err := git.RepositoryFiles.GetRawFile(args.Pid, fi.Path, &gitlab.GetRawFileOptions{
			Ref: branch,
		})

		if err != nil {
			log.Println(err)
			return
		}

		countSplit := strings.Split(fi.Path, "/")

		filename := ""
		outPath := ""

		outPath = strings.TrimSuffix(args.Out, "/")

		// 创建目录
		dir := createPath(outPath, countSplit, args)

		var build strings.Builder
		build.WriteString(dir)
		build.WriteString("/")
		build.WriteString(fi.Name)
		filename = build.String()

		generalWrite(filename, file)
	}

}
