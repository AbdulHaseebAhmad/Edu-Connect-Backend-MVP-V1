package Configurator

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address string `yaml:"Address" env:"Address" env-required:"true"`
}

type SMTP struct {
	Host     string `yaml:"Host" env:"Host" env-required:"true"`
	Port     string `yaml:"Port" env:"Port" env-required:"true"`
	Sender   string `yaml:"Sender" env:"Sender" env-required:"true"`
	Password string `yaml:"Password" env:"Password" env-required:"true"`
}

type Configuration struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server" env-required:"true"`
	SMTP        `yaml:"smtp" env-required:"true"`
}

func LoadConfiguration() *Configuration {

	var configPath string
	configPath = os.Getenv("../Config/localConfig.yaml") //using os package we call the Getenv method to look for environment variables in environment

	if configPath == "" {
		flags := flag.String("config_path", "Config/localConfig.yaml", "The path to the yaml file holding details")
		flag.Parse()
		configPath = *flags
		slog.Info("The path to the yaml holding environment variable:", "path", configPath)
		if configPath == "" {
			slog.Info("Add flags while running the command")
		}
	}

	_, error := os.Stat(configPath) //check if yaml file exist at the passed path from cli

	if os.IsNotExist(error) {
		slog.Info("config file does not exist: ", "path", configPath)
	}

	var Configuration Configuration
	err := cleanenv.ReadConfig(configPath, &Configuration)

	if err != nil {
		slog.Info("can not read config file: ", "error", err.Error())
	}

	return &Configuration

}
