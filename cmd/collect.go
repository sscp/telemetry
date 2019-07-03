package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sscp/telemetry/collector"

	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func registerCollectCmd(rootCmd *cobra.Command, rootConfig *viper.Viper) {
	// collectCmd represents the collect command
	var collectCmd = &cobra.Command{
		Use:   "collect",
		Short: "collects car data",
		Long:  `TODO: long-form doc`,
		Args:  cobra.ExactArgs(1),
		RunE:  createRunCollectorFunc(rootConfig),
	}
	// Flags ignored by UnmarshalKey...
	//collectCmd.PersistentFlags().IntP("port", "p", 33333, "port to listen for packets on")
	//rootConfig.BindPFlag("collector.port", collectCmd.PersistentFlags().Lookup("port"))

	rootCmd.AddCommand(collectCmd)
}

func createRunCollectorFunc(rootConfig *viper.Viper) func(cmd *cobra.Command, args []string) error {

	return func(cmd *cobra.Command, args []string) error {
		config := collector.CollectorConfig{}
		err := rootConfig.UnmarshalKey("collector", &config)
		if err != nil {
			return err
		}

		col, err := collector.NewUDPCollector(config)
		if err != nil {
			return err
		}

		fmt.Printf("Now collecting packets on port %v\n", config.Port)
		col.RecordRun(context.TODO(), args[0])
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("Press any key to end")
		_, err = buf.ReadBytes('\n')
		if err != nil {
			return err
		}
		return nil
	}

}
