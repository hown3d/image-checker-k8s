package k8s

import (
	"context"
	"errors"
	"sync"

	"github.com/hown3d/image-checker-k8s/pkg"
	"github.com/sirupsen/logrus"

	"github.com/containers/image/v5/types"
	//v1 "k8s.io/api/core/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesConfig struct {
	Namespaces       []string
	KubeConfig       string
	KubeClient       *kubernetes.Clientset
	RegistryOpts     *RegistryOption
	UpdateAnnotation string
}

type podOwnerMetaData struct {
	pod        *apiv1.Pod
	containers *[]apiv1.ContainerStatus
	ownerName  string
	kind       string
}

func (k *KubernetesConfig) NewClientSet() (err error) {
	// kubeconfig is set to the current set context in kubeConfig File
	//If neither masterUrl or kubeconfigPath are passed in we fallback to inClusterConfig. If inClusterConfig fails, we fallback to the default config.
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", k.KubeConfig)
	k.KubeClient, err = kubernetes.NewForConfig(kubeConfig)
	return
}
func (k *KubernetesConfig) getPodOwner(namespace string, pod *apiv1.Pod) (*podOwnerMetaData, error) {
	ctx := context.Background()
	getOpts := metav1.GetOptions{}
	appsClient := k.KubeClient.AppsV1()
	for _, owner := range pod.OwnerReferences {
		switch owner.Kind {
		case "ReplicaSet":
			replicaSet, err := appsClient.ReplicaSets(namespace).Get(ctx, owner.Name, getOpts)
			if err != nil {
				return nil, err
			}
			// dont know why this is a for loop, idk if a pod can have more then 1 owner
			for _, rsOwner := range replicaSet.OwnerReferences {
				return &podOwnerMetaData{pod, &pod.Status.ContainerStatuses, rsOwner.Name, rsOwner.Kind}, nil
			}
		case "DaemonSet":
			daemonSet, err := appsClient.DaemonSets(namespace).Get(ctx, owner.Name, getOpts)
			if err != nil {
				return nil, err
			}
			return &podOwnerMetaData{pod, &pod.Status.ContainerStatuses, daemonSet.Name, daemonSet.Kind}, nil
		case "StatefulSet":
			statefulSet, err := appsClient.StatefulSets(namespace).Get(ctx, owner.Name, getOpts)
			if err != nil {
				return nil, err
			}
			return &podOwnerMetaData{pod, &pod.Status.ContainerStatuses, statefulSet.Name, statefulSet.Name}, nil
		}
	}
	return nil, errors.New("Pod has no Owners, ???")

}

func (k *KubernetesConfig) updateTemplateSpec(podTemplate *apiv1.PodTemplateSpec, newImage string) {
	for index := range podTemplate.Spec.Containers {
		//Use index instead of second value, because that returns a copy, not a reference
		podTemplate.Spec.Containers[index].Image = newImage
	}
}

func (k *KubernetesConfig) updateResource(newImage string, podOwnerMeta *podOwnerMetaData, sysCtx *types.SystemContext) error {
	ctx := context.Background()
	getOpts := metav1.GetOptions{}
	updateOpts := metav1.UpdateOptions{}
	switch podOwnerMeta.kind {
	case "Deployment":
		deploymentClient := k.KubeClient.AppsV1().Deployments(podOwnerMeta.pod.Namespace)
		deployment, err := deploymentClient.Get(ctx, podOwnerMeta.ownerName, getOpts)
		logrus.Infof("Updating Deployment %v", deployment.Name)
		if err != nil {
			return err
		}
		k.updateTemplateSpec(&deployment.Spec.Template, newImage)
		_, exists := deployment.Annotations[k.UpdateAnnotation]
		if exists == false {
			deployment.Annotations[k.UpdateAnnotation] = "kekw"
		}
		_, err = deploymentClient.Update(ctx, deployment, updateOpts)
		return err

		//case "DaemonSet":

	}
	return nil
}

func (k *KubernetesConfig) GetRessourcesToUpdate(listOpts *metav1.ListOptions, sysCtx *types.SystemContext) error {

	for _, namespace := range k.Namespaces {
		logrus.Infof("Checking Namespace %v", namespace)
		pods, err := k.KubeClient.CoreV1().Pods(namespace).List(context.Background(), *listOpts)
		if err != nil {
			return err
		}
		for _, pod := range pods.Items {
			podOwnerData, err := k.getPodOwner(namespace, &pod)
			if err != nil {
				return err
			}
			logrus.Infof("Checking pod %v", pod.Name)
			k.updateSetIfNeeded(podOwnerData, sysCtx)
		}

	}
	return nil
}

func (k *KubernetesConfig) updateSetIfNeeded(podOwnerMeta *podOwnerMetaData, sysCtx *types.SystemContext) error {
	for _, currentContainer := range *podOwnerMeta.containers {
		if k.RegistryOpts.IsNewImage(currentContainer.Image, currentContainer.ImageID, sysCtx) {
			logrus.Infof("Container %v needs update, new digest found for image %v", currentContainer.Name, currentContainer.Image)
			_, registryDigest := k.RegistryOpts.GetDigests(currentContainer.Image, currentContainer.ImageID, sysCtx)
			newImage := currentContainer.Image + "@" + registryDigest
			logrus.Infof("New Image is %v", newImage)
			return k.updateResource(newImage, podOwnerMeta, sysCtx)
		}

	}
	return nil
}

func (k *KubernetesConfig) GetImageOfContainers(
	listOpts *metav1.ListOptions,
	tabWriter *pkg.TabWriter,
	sysCtx *types.SystemContext) (err error) {

	c := make(chan apiv1.ContainerStatus)
	var wg sync.WaitGroup
	wg.Add(10)
	for ii := 0; ii < 10; ii++ {
		go func(c chan apiv1.ContainerStatus) {
			for {
				currentContainer, more := <-c
				if more == false {
					wg.Done()
					return
				}
				needsChange := ""

				if k.RegistryOpts.IsNewImage(currentContainer.Image, currentContainer.ImageID, sysCtx) {
					needsChange = "X"
				}
				installedDigest, registryDigest := k.RegistryOpts.GetDigests(currentContainer.Image, currentContainer.ImageID, sysCtx)
				//print to stdout
				toPrint := []string{currentContainer.Image, installedDigest[:25], registryDigest[:25], needsChange}
				tabWriter.Write(toPrint...)
			}
		}(c)
	}
	for _, namespace := range k.Namespaces {
		pods, err := k.KubeClient.CoreV1().Pods(namespace).List(context.Background(), *listOpts)
		if err != nil {
			return err
		}
		for _, pod := range pods.Items {
			for _, container := range pod.Status.ContainerStatuses {
				c <- container
			}
		}

	}
	close(c)
	wg.Wait()
	return err
}
