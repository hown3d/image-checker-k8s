package cmd

import (
	"github.com/spf13/cobra"
)

func logoutCmd(opts *Options) *cobra.Command {

	// cmd represents the login command
	var cmd = &cobra.Command{
		Use:   "logout [REGISTRY]",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run:  opts.logout,
		Args: cobra.MaximumNArgs(1),
	}
	cmd.Flags().BoolVarP(&opts.K8s.RegistryOpts.LogOutFromAllRegistries, "all", "a", false, "specify if you want to logout from all registries")
	return cmd
}

func (opts *Options) logout(_ *cobra.Command, args []string) {
	registryName := args[0]
	opts.K8s.RegistryOpts.LogoutFromRegistry(registryName, opts.SysCtx)
}
