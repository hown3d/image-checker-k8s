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
}

func (opts *RegistryOption) GetDigests(image, imageID string, sysCtx *types.SystemContext) (installedDigest, registryDigest string) {
	installedDigest = strings.Split(imageID, "@")[1]
	registryDigest = opts.GetRegistryDigest(image, sysCtx).String()
	return
}

func (opts *RegistryOption) IsNewImage(image, imageID string, sysCtx *types.SystemContext) bool {
	installedDigest, registryDigest := opts.GetDigests(image, imageID, sysCtx)
	if installedDigest != registryDigest {
		return true
	}
	return false
}

func (opts *RegistryOption) LogoutFromRegistry(registryName string, sysCtx *types.SystemContext) error {
	logoutOpts := &auth.LogoutOptions{
		All:    opts.LogOutFromAllRegistries,
		Stdout: os.Stdout,
	}

	return auth.Logout(sysCtx, logoutOpts, []string{registryName})

}

func (opts *RegistryOption) LoginToRegistry(registryName string, sysCtx *types.SystemContext) error {
	loginOpts := &auth.LoginOptions{
		Username: opts.RegistryUser,
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
	}
	if opts.RegistryPassword != "" {
		loginOpts.Password = opts.RegistryPassword
	}
	return auth.Login(context.Background(), sysCtx, loginOpts, []string{registryName})
}

func (opts *RegistryOption) GetRegistryDigest(imageName string, sysCtx *types.SystemContext) *digest.Digest {

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

	image, err := reference.NewImage(ctx, sysCtx)
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
