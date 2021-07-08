
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)


func updateCmd(opts *Options) *cobra.Command{

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
	return cmd
}


func (opts *Options) update(_ *cobra.Command, args []string) {
	opts.createConfig()
	fmt.Println("not implemented")
}
