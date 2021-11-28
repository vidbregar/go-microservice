package config

type Config struct {
	Logger Logger `yaml:"logger"`
	Redis  Redis  `yaml:"redis"`
}

type Logger struct {
	Development bool `yaml:"development"`
}

type Redis struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	Retries int    `json:"retries"`
}
