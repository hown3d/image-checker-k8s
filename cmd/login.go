package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)



func loginCmd(opts *Options) *cobra.Command{

	// cmd represents the login command
	var cmd = &cobra.Command{
		Use:   "login [REGISTRY]",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: opts.login,
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().StringVarP(&opts.RegOptions.RegistryUser, "username", "u","", "username to login to registry (required if registry = true)")
	cmd.Flags().StringVarP(&opts.RegOptions.RegistryPassword, "password", "p", "", "password to login to registry")
	cmd.Flags().StringVar(&opts.RegOptions.AuthFile, "authfile", os.Getenv("REGISTRY_AUTH_FILE"), "path of the authentication file. Use REGISTRY_AUTH_FILE environment variable to override")
	return cmd
}


func (opts *Options) login(_ *cobra.Command, args []string) {
	opts.createConfig()
	registryName := args[0]

	err := opts.LoginToRegistry(registryName)
	if err != nil {
		log.Errorf("Can't log into registry %v, because %v", registryName, err)
	}
}
