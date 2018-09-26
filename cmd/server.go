package cmd

import (
	"fmt"
	"log"

	"github.com/sscp/telemetry/collector"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func registerServerCmd(rootCmd *cobra.Command, rootConfig *viper.Viper) {
	// serverCmd represents the server command
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Runs the telemetry server",
		Long:  `TODO: long-form doc`,
		Run:   createRunServerFunc(rootConfig),
	}
	// Bind doesn't work when unmarshaling...
	//serverCmd.PersistentFlags().IntP("port", "p", 3000, "port to listen on")
	//rootConfig.BindPFlag("server.port", serverCmd.PersistentFlags().Lookup("port"))

	rootCmd.AddCommand(serverCmd)

}

func createRunServerFunc(rootConfig *viper.Viper) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		config := collector.CollectorServiceConfig{}
		err := rootConfig.UnmarshalKey("server", &config)
		if err != nil {
			log.Fatalf("invalid config: %v", err)
		}

		fmt.Printf("Starting server on port %v, collector listening on port %v\n", config.Port, config.Collector.Port)
		collector.RunCollectionService(config)
	}

}
