package config

import (
	"flag"
	"path/filepath"

	"github.com/bahybintang/kube-restart/pkg/model"
	"github.com/bahybintang/kube-restart/pkg/utils"
	"github.com/caarlos0/env"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/util/homedir"
)

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" envDocs:"File path for kube-restart config" envDefault:"/etc/kube-restart.yml"`
}

type AppConfig struct {
	Deployments  []*model.App `yaml:"deployments" validate:"dive"`
	StatefulSets []*model.App `yaml:"statefulsets" validate:"dive"`
}

func (c *AppConfig) validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	return err
}

var (
	kubeconfig *string
)

func init() {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
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

func GetKubeConfig() *string {
	return kubeconfig
}
