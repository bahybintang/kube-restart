package service

import (
	"fmt"
	"time"

	"github.com/bahybintang/kube-restart/pkg/config"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

type UpdateFunc func(oldObj, newObj interface{})

func DeploymentWatcher(s *Store, cron *cron.Cron) (cache.Controller, error) {
	return Watcher(&schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}, func(oldObj, newObj interface{}) {
		u := newObj.(*unstructured.Unstructured)
		schedule := u.GetAnnotations()["kube.restart/schedule"]
		app := &StrippedApp{
			Name:      u.GetName(),
			Namespace: u.GetNamespace(),
		}
		if ValidateSchedule(schedule) && !s.IsDeploymentPresentWithTheSameConfig(app, schedule) {
			types, err := s.AddOrUpdateDeployment(app, schedule)
			if err != nil {
				logrus.Error("Failed to register deployment: ", err)
			}
			if types {
				logrus.Info(fmt.Sprintf("Added restart schedule (deployment): %v (%v) every %v", app.Name, app.Namespace, schedule))
			} else {
				logrus.Info(fmt.Sprintf("Updated restart schedule (deployment): %v (%v) every %v", app.Name, app.Namespace, schedule))
			}
		}
	})
}

func StatefulSetWatcher(s *Store, cron *cron.Cron) (cache.Controller, error) {
	return Watcher(&schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}, func(oldObj, newObj interface{}) {
		u := newObj.(*unstructured.Unstructured)
		schedule := u.GetAnnotations()["kube.restart/schedule"]
		app := &StrippedApp{
			Name:      u.GetName(),
			Namespace: u.GetNamespace(),
		}
		if ValidateSchedule(schedule) && !s.IsStatefulSetPresentWithTheSameConfig(app, schedule) {
			types, err := s.AddOrUpdateStatefulSet(app, schedule)
			if err != nil {
				logrus.Error("Failed to register statefulset: ", err)
			}
			if types {
				logrus.Info(fmt.Sprintf("Added restart schedule (statefulset): %v (%v) every %v", app.Name, app.Namespace, schedule))
			} else {
				logrus.Info(fmt.Sprintf("Updated restart schedule (statefulset): %v (%v) every %v", app.Name, app.Namespace, schedule))
			}
		}
	})
}

func Watcher(resource *schema.GroupVersionResource, updateFunc UpdateFunc) (cache.Controller, error) {
	client, err := config.GetKubeClient()
	if err != nil {
		return nil, err
	}
	dynInformer := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, time.Minute, corev1.NamespaceAll, nil)
	controller := dynInformer.ForResource(*resource).Informer()

	controller.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: updateFunc,
	})

	return controller, nil
}
