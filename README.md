# k8s yaml 资源生成工具

## linux平台构建

```bash
set CGO_ENABLED=0
set GOOS=linux
set GOPACH=amd64

go build -o deploymentool main.go
```

## 功能

- 生成dockerfile
- 下载gitlab中的文件
- 通过传递的json字符串，生成k8s yaml文件
- 通过简化的yaml文件，生成k8s yaml文件

## 使用

### 生成dockerfile

```bash
./deploymentool dockerfile 


# 参数
 args := struct {
  FromImg      string `cortana:"--img, -i, '', 基础镜像"`
  ProjectName  string `cortana:"--project, -p, '', 项目名称"`
  Outputs      string `cortana:"--outputs, -op, '', 制品路径"`
  JvmOpt       string `cortana:"--jvmOpt, -jvm,'', jvm参数"`
  Config       string `cortana:"--config, -c,'', 指定配置文件"`
  TemplateFile string `cortana:"--template, -t, '', dockerfile模板文件"`
  OutFile      string `cortana:"--out, -o, Dockerfile, 输出的文件名称"`
 }{}
```

### 下载gitlab中的文件

```bash
./deploymentool gitlabdown 

# 参数
 args := struct {
  Server    string `cortana:"--server, -s, , gitlab地址"`
  Branch    string `cortana:"--branch, -b, master, 分支"`
  ProjectId string `cortana:"--pid, -p, , gitlab项目id"`
  Path      string `cortana:"--path, , , gitlab项目id"`
  Encrypt   string `cortana:"--encrypt,  -e,'', 加密token字符串"`
  Token     string `cortana:"--token, -t, '',token加密后的字符串"`
  Out       string `cortana:"--out, -o, , 下载的路径"`
 }{}
```

### 通过传递的json字符串，生成k8s yaml文件

```
./deploymentool j2y 

# 参数
 args := struct {
  Context string `cortana:"--cfg, -c, , json字符串"`
  Write   string `cortana:"--write, -w, , 写入到文件"`
 }{}
```

### 通过简化的yaml文件，生成k8s yaml文件

```
./deploymentool deployfile 

# 参数
 args := struct {
  Config   string `cortana:"--cfg, -c, etc/kube-generator.yaml, 配置文件"`
  Json     string `cortana:"--json, -j, , 接收json字符串"`
  Template string `cortana:"--template, , template, 模板文件目录"`
  Tag      string `cortana:"--tag, -t, , 镜像标签"`
  Out      string `cortana:"--out, -o, , 输出的文件名称"`
  Dir      string `cortana:"--dir, -d, , 输出到目录"`
 }{}
```


## 配置文件

### default

```txt
Name         string
Namespace    string
Replicas     int
Labels       map[string]string
Image        string
Host         string
Ports        []PortConfig
Resources    map[string]ResourceConfig
HPA          []HPAConfig
Volumes      []map[interface{}]interface{}
VolumeMounts []map[interface{}]interface{}
```

default 为services的补充，类似于全局参数，services中未定义的参数，将使用default的参数

### services

```txt
Name         string
Namespace    string
Replicas     int
Labels       map[string]string
Image        string
Host         string
Ports        []PortConfig
Resources    map[string]ResourceConfig
HPA          []HPAConfig
Volumes      []map[interface{}]interface{}
VolumeMounts []map[interface{}]interface{}
```

#### name  名称
name  资源名称

#### namespace
name  kubernetes命名空间

#### replicas
replicas  pod副本数量

#### labels
labels 标签

#### image
image  镜像地址，可以通过 -t 传入

#### host
host    ingress域名,为空，不会生成ingres资源

#### Ports
ports 端口声明

#### Resources
resources   资源限制

#### HPA
hpa   自动扩容

#### Volumes
volumes  存储定义

#### VolumeMounts
volumemounts  存储挂载

## 示例

```yaml
default:
  namespace: dev
  replicas: 1
  resources:
    requests:
      cpu: 20m
      memory: 20Mi
    limits:
      cpu: 100m
      memory: 100Mi
      
  volumes:
    - name: test
      hostPath:
        path: /data
    - name: site-data
      persistentVolumeClaim:
        claimName: my-lamp-site-data

services:
  - name: be-kube-test-logdemo
    image: "reg.test.com/be-kube-test-logdemo:20220421103226_101"
    host: www.test.com
    ports:
      - port: 8080
        type: application
    hpa:
      - minreplicas: 1
        maxreplicas: 5
        type: cpu
        threshold: 80
      - minreplicas: 1
        maxreplicas: 5
        type: memory
        threshold: 80

  - name: be-kube-test-logdemo
    image: "reg.test.com/be-kube-test-logdemo:20220421103226_101"
    host: www.test.com
    ports:
      - port: 8080
        type: application
    volumemounts:
      - mountPath: /var/www/html
        name: site-data
        subPath: html
```
