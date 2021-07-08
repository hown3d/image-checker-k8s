package cmd

import (
	"context"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/hown3d/image-checker-k8s/k8s"
	"github.com/opencontainers/go-digest"

	"github.com/containers/image/v5/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

type Options struct {
	Ctx       context.Context
	SysCtx    *types.SystemContext
	TabWriter *tabwriter.Writer
	K8s       *k8s.KubernetesConfig
}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "image-checker-k8s",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Can't execute root command!")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	opts := Options{
		Ctx:       context.Background(),
		SysCtx:    &types.SystemContext{},
		TabWriter: new(tabwriter.Writer),
		K8s: &k8s.KubernetesConfig{
			RegistryOpts: &k8s.RegistryOption{
				DigestChache: make(map[string]*digest.Digest),
			},
		},
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	if home := homedir.HomeDir(); home != "" {
		rootCmd.PersistentFlags().StringVar(&opts.K8s.KubeConfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		rootCmd.PersistentFlags().StringVar(&opts.K8s.KubeConfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.image-checker-k8s.yaml)")

	//rootCmd.PersistentFlags().StringVar(&opts.RegOptions.AuthFile, "authfile", os.Getenv("REGISTRY_AUTH_FILE"), "path of the authentication file. Use REGISTRY_AUTH_FILE environment variable to override")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	//Add all subcommands
	rootCmd.AddCommand(
		loginCmd(&opts),
		updateCmd(&opts),
		listCmd(&opts),
		logoutCmd(&opts),
	)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			log.Errorln("Can't set homedir! ")
		}

		// Search config in home directory with name ".image-checker-k8s" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".image-checker-k8s")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %v", viper.ConfigFileUsed())
	}
}
