package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bahybintang/kube-restart/pkg/config"
	"github.com/bahybintang/kube-restart/pkg/model"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

func RestartDeployment(deployment *model.App) error {
	logrus.Info(fmt.Sprintf("Restarting deployment: %v (%v)", deployment.Name, deployment.Namespace))
	err := RestartApp(deployment, &schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"})
	if err != nil {
		logrus.Error(fmt.Sprintf("Deployment restart failed: %v (%v)", deployment.Name, deployment.Namespace))
	}
	logrus.Info(fmt.Sprintf("Deployment restarted: %v (%v)", deployment.Name, deployment.Namespace))
	return nil
}

func RestartStatefulSet(sts *model.App) error {
	logrus.Info(fmt.Sprintf("Restarting statefulset: %v (%v)", sts.Name, sts.Namespace))
	err := RestartApp(sts, &schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"})
	if err != nil {
		logrus.Error(fmt.Sprintf("Statefulset restart failed: %v (%v)", sts.Name, sts.Namespace))
	}
	logrus.Info(fmt.Sprintf("Statefulset restarted: %v (%v)", sts.Name, sts.Namespace))
	return nil
}

func RestartApp(app *model.App, resource *schema.GroupVersionResource) error {
	kubeconfig := config.GetKubeConfig()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err := client.Resource(*resource).Namespace(app.Namespace).Get(context.TODO(), app.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if err := unstructured.SetNestedField(result.Object, time.Now().UTC().Format("2006-01-02T15:04:05-0700"), "spec", "template", "metadata", "annotations", "kubectl.kubernetes.io/restartedAt"); err != nil {
			return err
		}
		_, err = client.Resource(*resource).Namespace(app.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return err
	})
}
