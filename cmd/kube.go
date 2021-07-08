package cmd

import (
	"context"
	"image-checker-k8s/pkg"
	"os"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/containers/common/pkg/auth"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/types"
	"github.com/opencontainers/go-digest"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	Ctx        context.Context
	SysCtx     *types.SystemContext
	KubeClient *kubernetes.Clientset
	TabWriter  *tabwriter.Writer
}
type registryOptions struct {
	RegistryPassword        string
	RegistryUser            string
	AuthFile                string
	LogOutFromAllRegistries bool
}

type Options struct {
	RegOptions   *registryOptions
	KubeConfig   string
	Config       *Config
	DigestChache map[string]*digest.Digest
	Namespaces   []string
}

func CreateClientSet(kubeConfigPath string) (clientset *kubernetes.Clientset, err error) {
	// kubeconfig is set to the current set context in kubeConfig File
	//If neither masterUrl or kubeconfigPath are passed in we fallback to inClusterConfig. If inClusterConfig fails, we fallback to the default config.
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	clientset, err = kubernetes.NewForConfig(kubeConfig)
	return
}

func (opts *Options) createConfig() {

	clientSet, err := CreateClientSet(opts.KubeConfig)
	if err != nil {
		log.Errorf("Can't create new kubernetes ClientSet")
	}

	opts.Config = &Config{
		Ctx: context.Background(),
		SysCtx: &types.SystemContext{
			AuthFilePath: opts.RegOptions.AuthFile,
		},
		KubeClient: clientSet,
		TabWriter:  new(tabwriter.Writer),
	}
}

func (opts *Options) GetImageOfContainers(namespaces []string) (podImages map[string]string, err error) {
	tabWriter := opts.Config.TabWriter
	podImages = make(map[string]string)

	err = pkg.TabWriterInit("Image\tInstalled Digest\tRegistry Digest", tabWriter)
	if err != nil {
		return nil, err
	}

	c := make(chan v1.ContainerStatus)
	var wg sync.WaitGroup
	wg.Add(10)
	for ii := 0; ii < 10; ii++ {
		go func(c chan v1.ContainerStatus) {
			for {
				currentContainer, more := <-c
				if more == false {
					wg.Done()
					return
				}

				//Name of Image (ex. docker.io/library/nginx:latest)
				imageID := currentContainer.Image

				installedDigest := strings.Split(currentContainer.ImageID, "@")[1]
				podImages[currentContainer.Name] = imageID
				registryDigest := opts.GetDigest(imageID)

				//print to stdout
				err = pkg.TabWriterWrite([]string{imageID, installedDigest[:25], registryDigest.String()[:25]}, tabWriter)
			}
		}(c)
	}
	for _, namespace := range namespaces {
		pods, err := opts.Config.KubeClient.CoreV1().Pods(namespace).List(opts.Config.Ctx, metav1.ListOptions{})
		if err != nil {
			return podImages, err
		}
		for _, pod := range pods.Items {
			for _, container := range pod.Status.ContainerStatuses {
				c <- container
			}
		}

	}
	close(c)
	wg.Wait()

	err = tabWriter.Flush()
	if err != nil {
		return nil, err

	}
	return podImages, err
}

func (opts *Options) LoginToRegistry(registryName string) error {
	loginOpts := &auth.LoginOptions{
		Username: opts.RegOptions.RegistryUser,
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
	}
	if opts.RegOptions.RegistryPassword != "" {
		loginOpts.Password = opts.RegOptions.RegistryPassword
	}
	return auth.Login(opts.Config.Ctx, opts.Config.SysCtx, loginOpts, []string{registryName})
}

func (opts *Options) GetDigest(imageName string) *digest.Digest {
	_, exists := opts.DigestChache[imageName]
	if exists == true {
		return opts.DigestChache[imageName]
	}

	reference, err := docker.ParseReference("//" + imageName)
	if err != nil {
		log.Fatalf("Can't parse reference on %v, because %v", imageName, err)
	}

	image, err := reference.NewImage(opts.Config.Ctx, opts.Config.SysCtx)
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

	manifestBytes, _, err := image.Manifest(opts.Config.Ctx)
	if err != nil {
		log.Errorf("Can't get manifest of image %v, because %v", imageName, err)
	}

	imageDigest, err := manifest.Digest(manifestBytes)
	if err != nil {
		log.Errorf("Can't get Digest of image %v, because %v", imageName, err)
	}

	opts.DigestChache[imageName] = &imageDigest

	return &imageDigest
}
