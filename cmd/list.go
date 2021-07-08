
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: listPods,
}

var namespaces []string

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().StringSliceVarP(&namespaces, "namespaces", "n", []string{"default"}, "namespaces to look for pods in")
}
func listPods(cmd *cobra.Command, args []string) {
	kubeConfigPath, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		log.Fatalf("Can't get kubeConfigPath, because %v", err)
	}
	config := createConfig(kubeConfigPath)

	_, err = config.GetImageOfContainers(namespaces)
	if err != nil {
		log.Errorf("Can't get Image of Containers, because %v", err)
	}
}