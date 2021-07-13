package cmd

import (
	"os"
	"text/tabwriter"

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
	cmd.Flags().BoolVarP(&opts.K8s.AllNamespaces, "all", "a", false, "use all namespaces for query")
	return cmd
}

func (opts *Options) list(_ *cobra.Command, _ []string) {

	if opts.K8s.KubeClient == nil {
		err := opts.K8s.NewClientSet()
		if err != nil {
			log.Fatal(err)
		}

	}

	if opts.K8s.Writer == nil {
		opts.K8s.Writer = os.Stdout
	}

	tabWriter := tabwriter.NewWriter(opts.K8s.Writer, 0, 8, 0, '\t', 0)

	_, err := tabWriter.Write([]byte("Image\tInstalled Digest\tRegistry Digest\tChange\n"))
	if err != nil {
		log.Fatal(err)
	}

	err = opts.K8s.GetImageOfContainers(&metav1.ListOptions{}, tabWriter)
	if err != nil {
		log.Errorf("Can't get Image of Containers, because %v", err)
	}

	err = tabWriter.Flush()
	if err != nil {
		log.Error(err)
	}

}
