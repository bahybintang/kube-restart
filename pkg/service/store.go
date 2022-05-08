package service

import (
	"github.com/bahybintang/kube-restart/pkg/config"
	"github.com/bahybintang/kube-restart/pkg/model"
	cron "github.com/robfig/cron/v3"
)

type StrippedApp struct {
	Name      string
	Namespace string
}

type Store struct {
	Cron                  *cron.Cron
	DeploymentsCronEntry  map[StrippedApp]cron.EntryID
	StatefulSetsCronEntry map[StrippedApp]cron.EntryID
	Deployments           map[StrippedApp]string
	StatefulSets          map[StrippedApp]string
}

func (s *Store) Init() {
	appConfig := config.GetAppConfig()
	for _, deploy := range appConfig.Deployments {
		deploy := deploy
		s.AddOrUpdateDeployment(&StrippedApp{
			Name:      deploy.Name,
			Namespace: deploy.Namespace,
		}, deploy.Schedule)
	}
	for _, statefulset := range appConfig.StatefulSets {
		statefulset := statefulset
		s.AddOrUpdateStatefulSet(&StrippedApp{
			Name:      statefulset.Name,
			Namespace: statefulset.Namespace,
		}, statefulset.Schedule)
	}
}

func (s *Store) AddOrUpdateDeployment(app *StrippedApp, schedule string) (bool, error) {
	_, ok := s.Deployments[*app]
	oldEntryId, ok1 := s.DeploymentsCronEntry[*app]

	entryId, err := s.Cron.AddFunc(schedule, func() {
		RestartDeployment(&model.App{
			Name:      app.Name,
			Namespace: app.Namespace,
			Schedule:  schedule,
		})
	})

	if err != nil {
		return false, err
	}

	if ok1 {
		s.Cron.Remove(oldEntryId)
	}

	s.Deployments[*app] = schedule
	s.DeploymentsCronEntry[*app] = entryId

	return !ok, nil
}

func (s *Store) AddOrUpdateStatefulSet(app *StrippedApp, schedule string) (bool, error) {
	_, ok := s.StatefulSets[*app]
	oldEntryId, ok1 := s.StatefulSetsCronEntry[*app]

	entryId, err := s.Cron.AddFunc(schedule, func() {
		RestartStatefulSet(&model.App{
			Name:      app.Name,
			Namespace: app.Namespace,
			Schedule:  schedule,
		})
	})

	if err != nil {
		return false, err
	}

	if ok1 {
		s.Cron.Remove(oldEntryId)
	}

	s.StatefulSets[*app] = schedule
	s.StatefulSetsCronEntry[*app] = entryId

	return !ok, nil
}
