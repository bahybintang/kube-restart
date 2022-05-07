package cmd

import (
	"github.com/bahybintang/kube-restart/pkg/config"
	"github.com/sirupsen/logrus"
)

func Main() {
	logrus.Info("Starting kube restart!")
	logrus.Info(config.GetAppConfig())
}
