package cmd

import (
	"time"

	"github.com/bahybintang/kube-restart/pkg/service"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func Main() {
	logrus.Info("Starting kube restart!")

	// Init store
	store := &service.Store{
		Cron:                  cron.New(),
		DeploymentsCronEntry:  make(map[service.StrippedApp]cron.EntryID),
		StatefulSetsCronEntry: make(map[service.StrippedApp]cron.EntryID),
		Deployments:           make(map[service.StrippedApp]string),
		StatefulSets:          make(map[service.StrippedApp]string),
	}
	store.Init()
	store.Cron.Start()

	// Start watcher
	deployController, err := service.DeploymentWatcher(store, store.Cron)
	if err != nil {
		logrus.Fatal("Failed to initialize deployment watcher: ", err)
	}
	deployStop := make(chan struct{})
	go deployController.Run(deployStop)

	stsController, err := service.StatefulSetWatcher(store, store.Cron)
	if err != nil {
		logrus.Fatal("Failed to initialize statefulset watcher: ", err)
	}
	stsStop := make(chan struct{})
	go stsController.Run(stsStop)

	for {
		time.Sleep(time.Second)
	}
}
