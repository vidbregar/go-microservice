package config

type Config struct {
	Server Server `yaml:"server"`
	Logger Logger `yaml:"logger"`
	Redis  Redis  `yaml:"redis"`
}

type Server struct {
	Address string `yaml:"address"`
}

type Logger struct {
	Development bool `yaml:"development"`
}

type Redis struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	Retries int    `json:"retries"`
}
