package config

import (
	"flag"

	"github.com/bahybintang/kube-restart/pkg/model"
	"github.com/bahybintang/kube-restart/pkg/utils"
	"github.com/caarlos0/env"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	kubeconfig = flag.String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
)

func init() {
	flag.Parse()
	logrus.Info(*kubeconfig)
	// Test kube config is valid
	_, err := GetKubeClient()
	if err != nil {
		logrus.Fatal("Failed to initialize kube client: ", err)
	}
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

func GetKubeClient() (dynamic.Interface, error) {
	var config *rest.Config
	var err error
	if *kubeconfig == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
