package cmd

import (
	"fmt"

	"github.com/sscp/telemetry/collector"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func registerServerCmd(rootCmd *cobra.Command, serverConfig *viper.Viper) {
	// serverCmd represents the server command
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Runs the telemetry server",
		Long:  `TODO: long-form doc`,
		Run:   createRunServerFunc(serverConfig),
	}
	// Flag causes crash?
	//serverCmd.PersistentFlags().IntP("port", "p", 3000, "port to listen on")
	//serverConfig.BindPFlag("port", serverCmd.Flags().Lookup("port"))

	rootCmd.AddCommand(serverCmd)
}

func createRunServerFunc(serverConfig *viper.Viper) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		config := collector.CollectorServiceConfig{}
		serverConfig.Unmarshal(&config)

		fmt.Printf("Starting server on port %v, collector listening on port %v\n", config.Port, config.Collector.Port)
		collector.RunCollectionService(config)
	}

}
