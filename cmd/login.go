package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func loginCmd(opts *Options) *cobra.Command {

	// cmd represents the login command
	var cmd = &cobra.Command{
		Use:   "login [REGISTRY]",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run:  opts.login,
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().StringVarP(&opts.K8s.RegistryOpts.RegistryUser, "username", "u", "", "username to login to registry (required if registry = true)")
	cmd.Flags().StringVarP(&opts.K8s.RegistryOpts.RegistryPassword, "password", "p", "", "password to login to registry")
	return cmd
}

func (opts *Options) login(_ *cobra.Command, args []string) {
	registryName := args[0]
	err := opts.K8s.RegistryOpts.LoginToRegistry(registryName)
	if err != nil {
		log.Errorf("%v", err)
	}
}
