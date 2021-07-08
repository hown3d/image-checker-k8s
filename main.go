package main

import (
	"image-checker-k8s/cmd"
)

func main() {
	//defer profile.Start(profile.TraceProfile, profile.ProfilePath("trace.out")).Stop()

	cmd.Execute()
}



