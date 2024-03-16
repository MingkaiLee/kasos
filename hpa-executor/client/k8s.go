package client

import (
	"context"
	"fmt"

	"github.com/MingkaiLee/kasos/hpa-executor/util"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var k8sClient *kubernetes.Clientset

const serviceNamespace = "default"

func InitK8sClient() {
	var err error
	conf, err := rest.InClusterConfig()
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
	k8sClient, err = kubernetes.NewForConfig(conf)
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
}

// 获取deployment
func GetDeployment(ctx context.Context, serviceName string) (d *v1.Deployment, err error) {
	deploymentName := fmt.Sprintf("%s-deployment", serviceName)
	d, err = k8sClient.AppsV1().Deployments(serviceNamespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		util.LogErrorf("client.GetDeployment error: %v", err)
	}
	return
}

// 更新deployment
func UpdateDeplotment(ctx context.Context, d *v1.Deployment) (err error) {
	_, err = k8sClient.AppsV1().Deployments(serviceNamespace).Update(ctx, d, metav1.UpdateOptions{})
	if err != nil {
		util.LogErrorf("client.UpdateDeplotment error: %v", err)
	}
	return
}
