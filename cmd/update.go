package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func updateCmd(opts *Options) *cobra.Command {

	// cmd represents the login command
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: opts.update,
	}
	cmd.Flags().StringSliceVarP(&opts.K8s.Namespaces, "namespaces", "n", []string{"default"}, "namespaces to look for pods in")
	cmd.Flags().BoolVarP(&opts.K8s.AllNamespaces, "all", "a", false, "if all namespaces should be checked for a new image")
	return cmd
}

func (opts *Options) update(_ *cobra.Command, args []string) {
	if opts.K8s.KubeClient == nil {
		err := opts.K8s.NewClientSet()
		if err != nil {
			log.Fatal(err)
		}

	}

	err := opts.K8s.GetRessourcesToUpdate(&metav1.ListOptions{})
	if err != nil {
		log.Error(err)
	}

}
