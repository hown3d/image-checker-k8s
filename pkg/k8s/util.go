package k8s

import (
	"context"
	"io"

	//v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesConfig struct {
	Namespaces       []string
	AllNamespaces    bool
	KubeConfig       string
	KubeClient       *kubernetes.Clientset
	RegistryOpts     *RegistryOption
	updateAnnotation string
	Writer           io.Writer
}

func (k *KubernetesConfig) NewClientSet() (err error) {
	//kubeconfig is set to the current set context in kubeConfig File
	//If neither masterUrl or kubeconfigPath are passed in we fallback to inClusterConfig. If inClusterConfig fails, we fallback to the default config.
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", k.KubeConfig)
	if err != nil {
		return
	}
	k.KubeClient, err = kubernetes.NewForConfig(kubeConfig)
	return
}

//Use Kubernetes API to return all Namespaces that are currently in the cluster
func (k *KubernetesConfig) getNamespaces() (namespaces []string, err error) {

	ctx := context.Background()
	if k.AllNamespaces {
		listOpts := metav1.ListOptions{}
		namespacesList, err := k.KubeClient.CoreV1().Namespaces().List(ctx, listOpts)
		if err != nil {
			return nil, err
		}
		for _, namespace := range namespacesList.Items {
			namespaces = append(namespaces, namespace.Name)
		}
	} else {
		namespaces = k.Namespaces
	}
	return namespaces, nil

}
