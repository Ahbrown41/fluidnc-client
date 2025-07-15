package cmd

import (
	"fmt"

	"fluidnc-client/internal/config"
	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/cobra"
)

var cmdCmd = &cobra.Command{
	Use:   "cmd [command]",
	Short: "Send single command via WebSocket",
	Long:  "Send a single command to FluidNC via WebSocket connection.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		client := fluidnc.NewClient(cfg)
		if err := client.Connect(); err != nil {
			return err
		}
		defer client.Disconnect()

		response, err := client.SendCommand(args[0])
		if err != nil {
			return err
		}

		fmt.Println(response)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cmdCmd)
}
