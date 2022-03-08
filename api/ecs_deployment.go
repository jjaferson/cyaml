package api

type ECSDeployment struct {
	Name     string    `yaml:"Name"`
	Network  Network   `yaml:"Network"`
	Services []Service `yaml:"Services"`
}

type Network struct {
	ID      string   `yaml:"Id"`
	Subnets []string `yaml:"Subnets"`
}

type Port struct {
	From int `yaml:"From"`
	To   int `yaml:"To"`
}

type Service struct {
	Name        string        `yaml:"Name"`
	Image       string        `yaml:"Image"`
	Replicas    int           `yaml:"Replicas"`
	Port        Port          `yaml:"Port"`
	CPU         string        `yaml:"Cpu"`
	Memory      string        `yaml:"Memory"`
	Environment []interface{} `yaml:"Environment"`
}
