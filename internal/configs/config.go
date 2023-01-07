package configs

import (
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
)

type Config struct {
	Server struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Timeout struct {
			Write time.Duration `yaml:"write"`
			Read  time.Duration `yaml:"read"`
		} `yaml:"timeout"`
	} `yaml:"server"`
	DB struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"DBName"`
		SSLMode  string `yaml:"SSLMode"`
	} `yaml:"db"`
}

func (cfg *Config) initENV() error {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	DB_PASSWORD, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		return errors.New("variable doesnt exists")
	}
	cfg.DB.Password = DB_PASSWORD
	return nil
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	err = config.initENV()
	if err != nil {
		return nil, err
	}
	return config, err
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags(path string) (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", path, "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}
