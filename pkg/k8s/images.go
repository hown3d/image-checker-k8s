package k8s

import (
	"context"
	"os"
	"strings"

	"github.com/containers/common/pkg/auth"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/manifest"

	"github.com/containers/image/v5/types"
	"github.com/opencontainers/go-digest"
	log "github.com/sirupsen/logrus"
)

type RegistryOption struct {
	RegistryPassword        string
	RegistryUser            string
	LogOutFromAllRegistries bool
	DigestChache            map[string]*digest.Digest
	SysCtx                  *types.SystemContext
}

func (opts *RegistryOption) GetDigests(image, imageID string) (installedDigest, registryDigest string) {
	log.Infof("Getting digests for image %v with ID %v", image, imageID)
	imageIDSplit := strings.Split(imageID, "@")

	//Case when imageID is only SHA256, not image@SHA256...
	if len(imageIDSplit) < 2 {
		installedDigest = imageIDSplit[0]
	} else {
		installedDigest = imageIDSplit[1]
	}

	registryDigest = opts.getRegistryDigest(image).String()
	return
}

func (opts *RegistryOption) IsNewImage(image, imageID string) bool {
	installedDigest, registryDigest := opts.GetDigests(image, imageID)
	if installedDigest != registryDigest {
		return true
	}
	return false
}

func (opts *RegistryOption) LogoutFromRegistry(registryName string) error {
	var logoutArgs []string
	logoutOpts := &auth.LogoutOptions{
		All:    opts.LogOutFromAllRegistries,
		Stdout: os.Stdout,
	}

	if logoutOpts.All == false {
		logoutArgs[0] = registryName
	}

	return auth.Logout(opts.SysCtx, logoutOpts, logoutArgs)

}

func (opts *RegistryOption) LoginToRegistry(registryName string) error {
	loginOpts := &auth.LoginOptions{
		Username: opts.RegistryUser,
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
	}
	if opts.RegistryPassword != "" {
		loginOpts.Password = opts.RegistryPassword
	}
	return auth.Login(context.Background(), opts.SysCtx, loginOpts, []string{registryName})
}

func (opts *RegistryOption) getRegistryDigest(imageName string) *digest.Digest {

	//Image can include digest (from earlier update)
	imageName = strings.Split(imageName, "@")[0]
	_, exists := opts.DigestChache[imageName]
	if exists == true {
		return opts.DigestChache[imageName]
	}

	reference, err := docker.ParseReference("//" + imageName)
	if err != nil {
		log.Fatalf("Can't parse reference on %v, because %v", imageName, err)
	}

	ctx := context.Background()

	image, err := reference.NewImage(ctx, opts.SysCtx)
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

	manifestBytes, _, err := image.Manifest(ctx)
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
