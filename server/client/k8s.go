package client

import (
	"context"
	"fmt"

	"github.com/MingkaiLee/kasos/server/util"
	monitorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var k8sClient *kubernetes.Clientset
var monitorClient *versioned.Clientset

const (
	monitorNamespace = "monitoring"
	serviceNamespace = "default"
)

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
	monitorClient, err = versioned.NewForConfig(conf)
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
}

// 获取ConfigMap
func GetConfigMap(ctx context.Context, name string) (configMap *corev1.ConfigMap, err error) {
	configMap, err = k8sClient.CoreV1().ConfigMaps(serviceNamespace).Get(ctx, name, metav1.GetOptions{})
	return
}

// 更新ConfigMap
func UpdateConfigMap(ctx context.Context, name string, configMap *corev1.ConfigMap) (err error) {
	_, err = k8sClient.CoreV1().ConfigMaps(serviceNamespace).Update(ctx, configMap, metav1.UpdateOptions{})
	return
}

func CreateMonitorService(ctx context.Context, serviceName string, tags map[string]string) (err error) {
	serviceMonitor := monitorv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-monitor", serviceName),
			Namespace: monitorNamespace,
		},
		Spec: monitorv1.ServiceMonitorSpec{
			Endpoints: []monitorv1.Endpoint{
				{
					Port: "http",
					Path: "/metrics",
				},
			},
			NamespaceSelector: monitorv1.NamespaceSelector{
				MatchNames: []string{
					serviceNamespace,
				},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: tags,
			},
		},
	}
	_, err = monitorClient.MonitoringV1().ServiceMonitors(monitorNamespace).Create(ctx, &serviceMonitor, metav1.CreateOptions{})
	return
}

func DeleteMonitorService(ctx context.Context, serviceName string) (err error) {
	err = monitorClient.MonitoringV1().ServiceMonitors(monitorNamespace).Delete(ctx, fmt.Sprintf("%s-monitor", serviceName), metav1.DeleteOptions{})
	return
}
