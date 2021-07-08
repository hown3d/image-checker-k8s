module image-checker-k8s

go 1.16

require (
	github.com/containers/image/v5 v5.13.2
	github.com/opencontainers/go-digest v1.0.0
	github.com/pkg/profile v1.6.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	k8s.io/api v0.20.6
	k8s.io/apimachinery v0.20.6
	k8s.io/client-go v0.20.6
)
