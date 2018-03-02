// Package cmd will form the basis of the telemetry command line interface. All
// commands should be unified in this package, which is called in the main.go
// of telemetry.
// TODO(jbeasley) This package is VERY primitive and currently only does the
// bare minimum to run something. I will flesh this out as we test collector
// and add more features/configurability
package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultCollectorConfig = map[string]interface{}{
		"port": 33333,
		"csv": map[string]interface{}{
			"folder": "./csvs",
		},
		"blog": map[string]interface{}{
			"folder": "./blogs",
		},
	}
	defaultServerConfig = map[string]interface{}{
		"port":      9090,
		"collector": defaultCollectorConfig,
	}
	defaultClientConfig = map[string]interface{}{
		"port":     9090,
		"hostname": "localhost",
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := buildRootTelemetryCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildRootTelemetryCmd() *cobra.Command {

	config := viper.New()
	config.SetDefault("client", defaultClientConfig)
	config.SetDefault("server", defaultServerConfig)
	config.SetDefault("collector", defaultCollectorConfig)

	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:               "telemetry",
		Short:             "telemetry manages the data that goes to and from the car",
		Long:              `telemetry manages the data that goes to and from the car`,
		PersistentPreRunE: createReadTelemetryConfigFunc(config),
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.telemetry.yaml)")

	registerServerCmd(rootCmd, config)
	registerCallCmd(rootCmd, config)
	registerCollectCmd(rootCmd, config)

	return rootCmd
}

func createReadTelemetryConfigFunc(config *viper.Viper) func(cmd *cobra.Command, args []string) error {

	return func(cmd *cobra.Command, args []string) error {

		cfgFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		if cfgFile != "" {
			// Use config file from the flag.
			config.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				return err
			}

			// Search config in home directory with name ".telemetry" (without extension).
			config.AddConfigPath(home)
			config.SetConfigName(".telemetry")
		}

		config.AutomaticEnv() // read in environment variables that match

		if config.ConfigFileUsed() != "" {
			// If a config file is found, read it in.
			if err := config.ReadInConfig(); err != nil {
				fmt.Println("Bad config")
				//return nil, err
			}
		}
		if config.ConfigFileUsed() != "" {
			fmt.Println("Using config file:", config.ConfigFileUsed())
		}

		return nil
	}
}
