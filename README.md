# Image-checker for Kubernetes
This tool tries to solve to the problem of containers running in kubernetes with an old build version of a certain tag. Without updating, securit updates might be missed and containers get vulnerable for attacks.

## Usage
```
$ image-checker-k8s --help

Usage:
  image-checker-k8s [command]

Available Commands:
  help        Help about any command
  list        list all containers and their images
  login       log in to remote container registry
  logout      logout from registry
  update      update containers with new images of remote registry

Flags:
      --config string       config file (default is $HOME/.image-checker-k8s.yaml)
  -h, --help                help for image-checker-k8s
      --kubeconfig string   (optional) absolute path to the kubeconfig file (default "/Users/hoehl/.kube/config")
  -t, --toggle              Help message for toggle

Use "image-checker-k8s [command] --help" for more information about a command.
```

## Development
- Build `go build`
- Run `./image-checker-k8s --help`
- Test `go test`

## Testing
Creates a image pull Secret in Cluster\
`kubectl create secret generic regcred \
--from-file=.dockerconfigjson=<path/to/.docker/config.json> \
--type=kubernetes.io/dockerconfigjson`
