package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: login,
	Args: cobra.ExactArgs(1),
}

var registryUser, registryPassword string
var privateRegistry bool

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")


	loginCmd.Flags().BoolVar(&privateRegistry, "private", false, "specify if registry is private or not")
	loginCmd.Flags().StringVarP(&registryUser, "username", "u","", "username to login to registry (required if registry = true)")
	loginCmd.Flags().StringVarP(&registryPassword, "password", "p", "", "password to login to registry")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func login(cmd *cobra.Command, args []string) {
	config := createConfig("")
	registryName := args[0]

	err := config.CheckAccessToRegistry(registryUser, registryPassword, registryName, privateRegistry)
	if err != nil {
		log.Errorf("Can't log into registry %v, because %v", registryName, err)
	} else {
		log.Infof("Successfully logged into registry %v", registryName)
	}
}
