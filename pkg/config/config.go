package config

type Config struct {
	Logger Logger `yaml:"logger"`
}

type Logger struct {
	Development bool `yaml:"development"`
}
