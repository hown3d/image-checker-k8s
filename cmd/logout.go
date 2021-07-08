package cmd

import (
	"github.com/containers/common/pkg/auth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)



func logoutCmd(opts *Options) *cobra.Command{

	// cmd represents the login command
	var cmd = &cobra.Command{
		Use:   "logout [REGISTRY]",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: opts.logout,
		Args: cobra.MaximumNArgs(1),

	}
	cmd.Flags().BoolVarP(&opts.RegOptions.LogOutFromAllRegistries, "all", "a", false,"specify if you want to logout from all registries")
	return cmd
}


func (opts *Options) logout(_ *cobra.Command, args []string) {
	opts.createConfig()

	logoutOpts := &auth.LogoutOptions{
		AuthFile: opts.RegOptions.AuthFile,
		All: opts.RegOptions.LogOutFromAllRegistries,
		Stdout: os.Stdout,
	}

	err := auth.Logout(opts.Config.SysCtx, logoutOpts, args)
	if err != nil {
		log.Errorf("Can't logout from registry, because %v", err)
	}
}
