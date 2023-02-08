package config

//import "github.com/flosch/pongo2/v6"

type DockerfileTemplate struct {
	Context map[string]any
	TplFile string
	OutFile string
}
