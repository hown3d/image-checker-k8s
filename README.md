# Image-checker for Kubernetes

## Create image pull Secret
- Path of dockerconfig is specified with `image-checker-k8s login <REGISTRY> --authfile`


## Testing
Creates a image pull Secret in Cluster
`kubectl create secret generic regcred \
--from-file=.dockerconfigjson=<path/to/.docker/config.json> \
--type=kubernetes.io/dockerconfigjson`
