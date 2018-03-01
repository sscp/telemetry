package cmd

import (
	"fmt"
	"log"

	"github.com/sscp/telemetry/collector"
	pb "github.com/sscp/telemetry/collector/serviceproto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func readConfAndConnect(rootConfig *viper.Viper) *collector.CollectorClient {
	config := collector.CollectorClientConfig{}
	rootConfig.Unmarshal(&config)

	client, err := collector.NewCollectorClient(config)
	if err != nil {
		log.Fatalf("Could not connect to collector: %v", err)
	}
	return client

}

func printCollectorStatus(status *pb.CollectorStatus) {
	if status.Collecting {
		fmt.Printf("Collector status: collecting run %v from port %v. %v packets received\n", status.RunName, status.Port, status.PacketsRecorded)
	} else {
		fmt.Print("Collector status: not collecting")
	}
}

func registerCallCmd(rootCmd *cobra.Command, clientConfig *viper.Viper) {
	// callCmd represents the call command
	var callCmd = &cobra.Command{
		Use:   "call",
		Short: "calls a collector server using GRPC",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			client := readConfAndConnect(clientConfig)
			defer client.Close()
		},
	}
	//callCmd.PersistentFlags().IntP("port", "p", 3000, "port to connect to")
	//clientConfig.BindPFlag("port", callCmd.Flags().Lookup("port"))

	//callCmd.PersistentFlags().StringP("host", "h", "localhost", "host to connect to")
	//clientConfig.BindPFlag("host", callCmd.Flags().Lookup("host"))

	// startCmd represents the start command
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "start a collector server over GRPC",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run:   createCallStart(clientConfig),
	}
	callCmd.AddCommand(startCmd)

	// stopCmd represents the stop command
	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "stop a collector server over GRPC",
		Long:  ``,
		Run:   createCallStop(clientConfig),
	}
	callCmd.AddCommand(stopCmd)

	// statusCmd represents the status command
	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "get the status of a collector server over GRPC",
		Long:  ``,
		Run:   createCallStatus(clientConfig),
	}
	callCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(callCmd)
}

func createCallStart(clientConfig *viper.Viper) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		runName := args[0]
		config := collector.CollectorClientConfig{}
		clientConfig.Unmarshal(&config)

		client := readConfAndConnect(clientConfig)
		defer client.Close()

		fmt.Printf("Starting collection for %v on port %v\n", config.Hostname, runName)
		status, err := client.StartCollector(runName)
		if err != nil {
			log.Fatalf("Could not start collector: %v", err)
		}
		fmt.Println("Started collector")
		printCollectorStatus(status)
	}
}

func createCallStop(clientConfig *viper.Viper) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		config := collector.CollectorClientConfig{}
		clientConfig.Unmarshal(&config)

		client := readConfAndConnect(clientConfig)
		defer client.Close()
		fmt.Println("Stopping collection")
		status, err := client.StopCollector()
		if err != nil {
			log.Fatalf("Could not stop collector: %v", err)
		}
		fmt.Println("Stopped collector")
		printCollectorStatus(status)
	}
}

func createCallStatus(clientConfig *viper.Viper) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		config := collector.CollectorClientConfig{}
		clientConfig.Unmarshal(&config)

		client := readConfAndConnect(clientConfig)
		defer client.Close()

		status, err := client.GetCollectorStatus()
		if err != nil {
			log.Fatalf("Could not get status: %v", err)
		}
		printCollectorStatus(status)
	}
}
