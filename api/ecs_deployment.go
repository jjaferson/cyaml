package api

type ECSDeployment struct {
	Name     string     `yaml:"Name"`
	Network  Network    `yaml:"Network"`
	Services []Services `yaml:"Services"`
}

type Network struct {
	ID      string   `yaml:"Id"`
	Subnets []string `yaml:"Subnets"`
}

type Port struct {
	From int `yaml:"From"`
	To   int `yaml:"To"`
}

type Services struct {
	Name        string        `yaml:"Name"`
	Image       string        `yaml:"Image"`
	Replicas    int           `yaml:"Replicas"`
	Port        Port          `yaml:"Port"`
	CPU         int           `yaml:"Cpu"`
	Memory      int           `yaml:"Memory"`
	Environment []interface{} `yaml:"Environment"`
}
