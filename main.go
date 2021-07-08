package main

import (
	"github.com/hown3d/image-checker-k8s/cmd"
)

func main() {
	//defer profile.Start(profile.TraceProfile, profile.ProfilePath("trace.out")).Stop()

	cmd.Execute()
}
