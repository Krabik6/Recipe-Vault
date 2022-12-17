package apiserver

// Config ...
type Config struct {
	Port string `yaml:"port"`
	//LogLevel string `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		Port: ":8000",
	}
}
