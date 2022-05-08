package cmd

import (
	"github.com/bahybintang/kube-restart/pkg/config"
	"github.com/bahybintang/kube-restart/pkg/service"
	"github.com/sirupsen/logrus"
)

func Main() {
	logrus.Info("Starting kube restart!")
	appConfig := config.GetAppConfig()
	logrus.Info(appConfig)
	service.RestartDeployment(appConfig.Deployments[0])
	service.RestartStatefulSet(appConfig.StatefulSets[0])
}
