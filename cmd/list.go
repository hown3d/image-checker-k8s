
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)



func listCmd(opts *Options) *cobra.Command{

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
	cmd.Flags().StringSliceVarP(&opts.Namespaces, "namespaces", "n", []string{"default"}, "namespaces to look for pods in")
		return cmd
}


func (opts *Options) list(_ *cobra.Command, _ []string) {
	opts.createConfig()
	_, err := opts.GetImageOfContainers(opts.Namespaces)
	if err != nil {
		log.Errorf("Can't get Image of Containers, because %v", err)
	}
}
