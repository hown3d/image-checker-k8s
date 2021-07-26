package k8s

import (
	"context"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type containerImage struct {
	image   string
	imageid string
}

func (k *KubernetesConfig) GetImageOfContainers(
	listOpts *metav1.ListOptions,
	writer io.Writer) (err error) {

	c := make(chan containerImage)
	var wg sync.WaitGroup
	wg.Add(10)

	for ii := 0; ii < 10; ii++ {
		go func() {
			for {
				currentContainerImage, more := <-c
				if !more {
					wg.Done()
					return
				}
				needsChange := ""

				if k.RegistryOpts.IsNewImage(currentContainerImage.image, currentContainerImage.imageid) {
					needsChange = "X"
				}
				installedDigest, registryDigest := k.RegistryOpts.GetDigests(currentContainerImage.image, currentContainerImage.imageid)
				//print to stdout
				toPrint := []byte(currentContainerImage.image + "\t" + installedDigest[:25] + "\t" + registryDigest[:25] + "\t" + needsChange + "\n")
				writer.Write(toPrint)

			}
		}()
	}

	namespaces, err := k.getNamespaces()
	if err != nil {
		return
	}

	log.Infof("Listing all pods in Namespaces %v", namespaces)

	for _, namespace := range namespaces {
		pods, err := k.KubeClient.CoreV1().Pods(namespace).List(context.Background(), *listOpts)
		if err != nil {
			return err
		}
		for _, pod := range pods.Items {
			containerStatus := pod.Status.ContainerStatuses
			containers := pod.Spec.Containers
			for index := range containers {
				c <- containerImage{containers[index].Image, containerStatus[index].ImageID}
			}
		}

	}
	close(c)
	wg.Wait()
	return err
}
