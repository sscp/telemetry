package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sscp/telemetry/collector/sources"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func registerUDPSendCmd(rootCmd *cobra.Command, rootConfig *viper.Viper) {
	// udpsendCmd represents the collect command
	var udpsendCmd = &cobra.Command{
		Use:   "udpsend",
		Short: "sends test udp packets",
		Long:  `sends test udp packets`,
		Args:  cobra.ExactArgs(0),
		Run:   createUDPSendFunc(rootConfig),
	}
	// Flags ignored by UnmarshalKey...
	//collectCmd.PersistentFlags().IntP("port", "p", 33333, "port to listen for packets on")
	//rootConfig.BindPFlag("collector.port", collectCmd.PersistentFlags().Lookup("port"))

	rootCmd.AddCommand(udpsendCmd)
}

func createUDPSendFunc(rootConfig *viper.Viper) func(cmd *cobra.Command, args []string) {

	return func(cmd *cobra.Command, args []string) {
		zps := sources.NewZeroPacketSource(20)
		go sources.SendPacketsAsUDP(zps.Packets(), 33333)
		zps.Listen()
		fmt.Printf("Now sending packets on port %v\n", 33333)
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("Press any key to end")
		buf.ReadBytes('\n')
		zps.Close()

	}

}
