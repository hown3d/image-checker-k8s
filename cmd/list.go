package cmd

import (
	"github.com/hown3d/image-checker-k8s/pkg"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func listCmd(opts *Options) *cobra.Command {

	// cmd represents the login command
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: opts.list,
	}
	cmd.Flags().StringSliceVarP(&opts.K8s.Namespaces, "namespaces", "n", []string{"default"}, "namespaces to look for pods in")
	return cmd
}

func (opts *Options) list(_ *cobra.Command, _ []string) {

	if opts.K8s.KubeClient == nil {
		opts.K8s.NewClientSet()
	}

	tabWriter := opts.TabWriter
	err := pkg.TabWriterInit("Image\tInstalled Digest\tRegistry Digest\tChange", tabWriter)
	if err != nil {
		log.Errorf("Can't write to tabWriter, because %v", err)
	}

	err = opts.K8s.GetImageOfContainers(&metav1.ListOptions{}, tabWriter, opts.SysCtx)
	if err != nil {
		log.Errorf("Can't get Image of Containers, because %v", err)
	}

	tabWriter.Flush()
}
