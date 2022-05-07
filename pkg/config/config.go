package config

import (
	"github.com/bahybintang/kube-restart/pkg/utils"
	"github.com/caarlos0/env"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" envDocs:"File path for kube-restart config" envDefault:"/etc/kube-restart.yml"`
}

type AppConfig struct {
	Deployments []struct {
		Name      string `yaml:"name" validate:"required"`
		Namespace string `yaml:"namespace" validate:"required"`
		Schedule  string `yaml:"schedule" validate:"required"`
	} `yaml:"deployments" validate:"dive"`

	Statefulsets []struct {
		Name      string `yaml:"name" validate:"required"`
		Namespace string `yaml:"namespace" validate:"required"`
		Schedule  string `yaml:"schedule" validate:"required"`
	} `yaml:"statefulsets" validate:"dive"`
}

func (c *AppConfig) validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	return err
}

func GetAppConfig() *AppConfig {
	appConfig := &AppConfig{}
	config := GetConfig()
	yamlFile, err := utils.OpenFile(config.ConfigPath)
	if err != nil {
		logrus.Fatal("Unable to read config file: ", err)
	}
	err = yaml.Unmarshal(yamlFile, appConfig)
	if err != nil {
		logrus.Fatal("Unable to unmarshal config file: ", err)
	}
	err = appConfig.validate()
	if err != nil {
		logrus.Fatal("Unable to validate config file: ", err)
	}
	return appConfig
}

func GetConfig() *Config {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		logrus.Fatal("Unable to parse environment variables: ", err)
	}
	return config
}
