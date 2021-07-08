package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/types"
	"github.com/opencontainers/go-digest"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"text/tabwriter"
)

type Config struct {
	Ctx        context.Context
	SysCtx     *types.SystemContext
	KubeClient *kubernetes.Clientset
	TabWriter  *tabwriter.Writer
	DigestChache map[string]*digest.Digest
}

func CreateClientSet(kubeConfigPath string) (clientset *kubernetes.Clientset, err error) {
	// kubeconfig is set to the current set context in kubeConfig File
	//If neither masterUrl or kubeconfigPath are passed in we fallback to inClusterConfig. If inClusterConfig fails, we fallback to the default config.
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	clientset, err = kubernetes.NewForConfig(kubeConfig)
	return
}

func (c Config) GetImageOfContainers(namespaces []string) (map[string]string, error) {
	c.TabWriter.Init(os.Stdout, 0, 8, 0, '\t', 0)
	_, err := fmt.Fprintln(c.TabWriter, "Pod\tImage\tDigest\tNamespace")
	if err != nil {
		return nil, err
	}
	podImages := make(map[string]string)



	for _, namespace := range namespaces {
		pods, err := c.KubeClient.CoreV1().Pods(namespace).List(c.Ctx, metav1.ListOptions{})
		if err != nil {
			return podImages, err
		}
		for _, pod := range pods.Items {
			for _, container := range pod.Status.ContainerStatuses {
					//imageSlice := strings.Split(container.ImageID, "@")
					imageID := container.Image

					podImages[container.Name] = imageID
					imageDigest := c.GetDigest(imageID)
					fmt.Fprintf(c.TabWriter, "%v\t%v\t%v\t%v\t\n", pod.Name, imageID, imageDigest.String(), namespace)
				}
			}


	}

	err = c.TabWriter.Flush()
	if err != nil {
		return nil, err

	}
	return podImages, nil
}
func (c *Config) CheckAccessToRegistry(username string, password string, registryName string, private bool) error {

	if (password == "" || username == "") && (private == true) {
		return errors.New("can't have empty user or password when registry is private")
	}
	return docker.CheckAuth(c.Ctx, c.SysCtx, username, password, registryName)

}

func (c *Config) GetDigest(imageName string) *digest.Digest {
	_, exists := c.DigestChache[imageName]
	if exists == true {
		return c.DigestChache[imageName]
	}

	reference, err := docker.ParseReference("//" + imageName)
	if err != nil {
		log.Fatalf("Can't parse reference on %v, because %v", imageName, err)
	}

	image, err := reference.NewImage(c.Ctx, c.SysCtx)
	if err != nil {
		log.Fatalf("Can't create new Image, because %v", err)
	}

	// Close Image on exit
	defer func(image types.ImageCloser) {
		err := image.Close()
		if err != nil {
			log.Errorf("Can't close Image %v, because %v", image, err)
		}
	}(image)

	manifestBytes, _, err := image.Manifest(c.Ctx)
	if err != nil {
		log.Errorf("Can't get manifest of image %v, because %v", imageName, err)
	}

	imageDigest, err := manifest.Digest(manifestBytes)
	if err != nil {
		log.Errorf("Can't get Digest of image %v, because %v", imageName, err)
	}

	c.DigestChache[imageName] = &imageDigest

	return &imageDigest
}
