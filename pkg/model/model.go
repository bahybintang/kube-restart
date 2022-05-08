package model

type App struct {
	Name      string `yaml:"name" validate:"required"`
	Namespace string `yaml:"namespace" validate:"required"`
	Schedule  string `yaml:"schedule" validate:"required"`
}
