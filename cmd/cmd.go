package cmd

import (
	"github.com/bahybintang/kube-restart/pkg/config"
	"github.com/bahybintang/kube-restart/pkg/service"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func Main() {
	logrus.Info("Starting kube restart!")
	appConfig := config.GetAppConfig()
	c := cron.New()
	for _, deployment := range appConfig.Deployments {
		c.AddFunc(deployment.Schedule, func() { service.RestartDeployment(deployment) })
	}
	for _, statefulset := range appConfig.StatefulSets {
		c.AddFunc(statefulset.Schedule, func() { service.RestartStatefulSet(statefulset) })
	}
	c.Start()
	for {
	}
}
