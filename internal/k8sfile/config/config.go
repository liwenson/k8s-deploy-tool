package config

type PortConfig struct {
	Port       int
	Type       string
	TargetPort int
}

type ResourceConfig struct {
	CPU    string
	Memory string
}

type HPAConfig struct {
	MinReplicas int
	MaxReplicas int
	Type        string
	Threshold   int
}

type ContentConfig struct {
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
}

type Config struct {
	OutputPath   string
	TemplatePath string
	Services     []ContentConfig
	Default      ContentConfig
}

type Args struct {
	ConfigFile string
	JsonStr    string
	Template   string
	Tag        string
	Out        string
	Dir        string
}
