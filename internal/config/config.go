package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env        string        `yaml:"env" env-default:"prod"`
	Storage    Storage       `yaml:"storage"`
	GRPCServer GRPCServer    `yaml:"grpc_server"`
	Clients    ClientsConfig `yaml:"clients"`
}

type Storage struct {
	Type string `yaml:"type" env-default:"memory"`
	Path string `yaml:"storage_path" env-default:"./storage/preview.db"`
}

type GRPCServer struct {
	Port int `yaml:"port"`
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
	OutputDir    string        `yaml:"output_dir"`
}

type ClientsConfig struct {
	Preview Client `yaml:"preview"`
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
