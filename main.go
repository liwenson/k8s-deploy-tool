package main

import (
	dockerfile_cfg "NBServer/tools/deploymentgen/internal/dockerfile/config"
	dockerfile_logic "NBServer/tools/deploymentgen/internal/dockerfile/logic"
	gitlabfile_cfg "NBServer/tools/deploymentgen/internal/getgitlabfile/config"
	gitlabfile_logic "NBServer/tools/deploymentgen/internal/getgitlabfile/logic"
	jsonToyaml_cfg "NBServer/tools/deploymentgen/internal/jsonToyaml/config"
	jsonToyaml_logic "NBServer/tools/deploymentgen/internal/jsonToyaml/logic"
	k8sfile_cfg "NBServer/tools/deploymentgen/internal/k8sfile/config"
	k8sfile_logic "NBServer/tools/deploymentgen/internal/k8sfile/logic"
	"fmt"
	"github.com/shafreeck/cortana"
	"strings"
)

////go:embed template/*
//var TempFs embed.FS

func jsonToYaml() {
	cortana.Title("json转yaml")
	cortana.Description(`将json字符串转换为yaml格式输出到控制台或文件`)

	args := struct {
		Context string `cortana:"--cfg, -c, , json字符串"`
		Write   string `cortana:"--write, -w, , 写入到文件"`
	}{}

	cortana.Parse(&args)

	fmt.Println(args.Context)

	conf := jsonToyaml_cfg.Content{
		Content: args.Context,
		Write:   args.Write,
	}

	jsonToyaml_logic.JsonToYaml(conf)
}

func deployFileGen() {
	cortana.Title("生成k8s deployment 资源")
	cortana.Description(`下载gitlab仓库中的文件或者目录
--out 路径带 / 将不保留第一层目录
`)

	args := struct {
		Config   string `cortana:"--cfg, -c, etc/kube-generator.yaml, 配置文件"`
		Json     string `cortana:"--json, -j, , 接收json字符串"`
		Template string `cortana:"--template, , template, 模板文件目录"`
		Tag      string `cortana:"--tag, -t, , 镜像标签"`
		Out      string `cortana:"--out, -o, , 输出的文件名称"`
		Dir      string `cortana:"--dir, -d, , 输出到目录"`
	}{}

	cortana.Parse(&args)

	conf := k8sfile_cfg.Args{
		ConfigFile: args.Config,
		JsonStr:    args.Json,
		Template:   args.Template,
		Tag:        args.Tag,
		Out:        args.Out,
		Dir:        args.Dir,
	}

	c := k8sfile_cfg.MustLoad(&conf)
	k8sfile_logic.GenerateServicesYaml(&c)
}

func gitlabFileGet() {
	cortana.Title("gitlab下载工具")
	cortana.Description(`下载gitlab仓库中的文件或者目录
--out 路径带 / 将不保留第一层目录
`)

	args := struct {
		Server    string `cortana:"--server, -s, , gitlab地址"`
		Branch    string `cortana:"--branch, -b, master, 分支"`
		ProjectId string `cortana:"--pid, -p, , gitlab项目id"`
		Path      string `cortana:"--path, , , gitlab项目id"`
		Encrypt   string `cortana:"--encrypt,  -e,'', 加密token字符串"`
		Token     string `cortana:"--token, -t, '',token加密后的字符串"`
		Out       string `cortana:"--out, -o, , 下载的路径"`
	}{}

	cortana.Parse(&args)

	if args.Encrypt != "" {
		gitlabfile_logic.PwdEncryption(args.Encrypt)
		return
	}

	conf := gitlabfile_cfg.Args{
		Server: args.Server,
		Branch: args.Branch,
		Path:   args.Path,
		Pid:    args.ProjectId,
		Out:    args.Out,
		Token:  gitlabfile_logic.PwdDecrypt(args.Token),
	}

	//vTtVUdGD6o2W5N5nTJDZALIW8OejBH92yff6RL0Kcok=

	// 检查参数是否为空
	gitlabfile_cfg.IsArgeNull(conf)
	gitlabfile_logic.GitlabRun(conf)
}

func dockerFileGen() {
	cortana.Title("生成dockerfile")
	cortana.Description(`通过模板文件生成dockerfile,
模板语法为Django模板语法`)

	args := struct {
		FromImg      string `cortana:"--img, -i, '', 基础镜像"`
		ProjectName  string `cortana:"--project, -p, '', 项目名称"`
		Outputs      string `cortana:"--outputs, -op, '', 制品路径"`
		JvmOpt       string `cortana:"--jvmOpt, -jvm,'', jvm参数"`
		Config       string `cortana:"--config, -c,'', 指定配置文件"`
		TemplateFile string `cortana:"--template, -t, '', dockerfile模板文件"`
		OutFile      string `cortana:"--out, -o, Dockerfile, 输出的文件名称"`
	}{}

	cortana.Parse(&args)

	var (
		cfg    []string
		scf    string
		dcf    string
		anchor int
	)

	if args.Config != "" {
		// 配置文件路径判断
		cfg = strings.Split(args.Config, ":")
	}

	if len(cfg) == 1 {
		scf = cfg[0]
		tcfg := strings.Split(cfg[0], "/")
		dcf = tcfg[len(tcfg)-1]
		anchor = 1
	} else if len(cfg) == 2 {
		scf = cfg[0]
		dcf = cfg[1]
		anchor = 2
	} else {
		anchor = 0
		//log.Fatalf("Failed to create client: %s", "参数不正确")
	}

	//jvmOpt := strings.Trim("@", args.JvmOpt)

	context := map[string]any{
		"FromImg":     args.FromImg,
		"ProjectName": args.ProjectName,
		"Outputs":     args.Outputs,
		"JvmOpt":      args.JvmOpt,
		//"Config":      args.Config,
		"Scf":    scf,
		"Dcf":    dcf,
		"Anchor": anchor, // 标记配置文件路径
	}

	dockerfile := dockerfile_cfg.DockerfileTemplate{
		Context: context,
		TplFile: args.TemplateFile,
		OutFile: args.OutFile,
		//TplFile: pongo2.Must(pongo2.FromFile("template/dockerfile.template")),
	}

	dockerfile_logic.Gendockerfile(dockerfile)

}

func main() {
	cortana.AddCommand("dockerfile", dockerFileGen, "print anything to the screen")
	cortana.AddCommand("gitlabdown", gitlabFileGet, "echo anything to the screen")
	cortana.AddCommand("deployfile", deployFileGen, "echo anything to the screen more times")
	cortana.AddCommand("j2y", jsonToYaml, "echo anything to the screen more times")
	cortana.Launch()
}
