package service

import (
	"fmt"
	"regexp"

	"github.com/bahybintang/kube-restart/pkg/config"
	"github.com/bahybintang/kube-restart/pkg/model"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
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

func (s *Store) IsDeploymentPresentWithTheSameConfig(app *StrippedApp, schedule string) bool {
	val, ok := s.Deployments[*app]
	logrus.Debug(fmt.Sprintf("Validating deployment: %v %v %v", app, val, ok))
	return ok && val == schedule
}

func (s *Store) IsStatefulSetPresentWithTheSameConfig(app *StrippedApp, schedule string) bool {
	val, ok := s.StatefulSets[*app]
	logrus.Debug(fmt.Sprintf("Validating sts: %v %v %v", app, val, ok))
	return ok && val == schedule
}

func ValidateSchedule(schedule string) bool {
	regex, _ := regexp.Compile(`(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})`)
	ret := regex.MatchString(schedule)
	logrus.Debug(fmt.Sprintf("Validating schedule: %v %v", schedule, ret))
	return ret
}
