package cmd

import (
	"fmt"

	"fluidnc-client/internal/config"
	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/cobra"
)

var httpControlCmd = &cobra.Command{
	Use:   "http-control",
	Short: "HTTP-based machine control commands",
	Long:  "Send machine control commands via HTTP endpoints instead of WebSocket.",
}

var httpHoldCmd = &cobra.Command{
	Use:   "hold",
	Short: "Send feed hold via HTTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		client := fluidnc.NewClient(cfg)
		if err := client.HTTPFeedHold(); err != nil {
			return err
		}
		fmt.Println("Feed hold sent via HTTP")
		return nil
	},
}

var httpStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Send cycle start via HTTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		client := fluidnc.NewClient(cfg)
		if err := client.HTTPCycleStart(); err != nil {
			return err
		}
		fmt.Println("Cycle start sent via HTTP")
		return nil
	},
}

var httpRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart FluidNC via HTTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		client := fluidnc.NewClient(cfg)
		if err := client.HTTPRestart(); err != nil {
			return err
		}
		fmt.Println("Restart command sent via HTTP")
		return nil
	},
}

var checkRestartCmd = &cobra.Command{
	Use:   "check-restart",
	Short: "Check if restart occurred",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		client := fluidnc.NewClient(cfg)
		restarted, err := client.CheckDidRestart()
		if err != nil {
			return err
		}
		if restarted {
			fmt.Println("Device has restarted")
		} else {
			fmt.Println("Device has not restarted")
		}
		return nil
	},
}

func init() {
	httpControlCmd.AddCommand(httpHoldCmd, httpStartCmd, httpRestartCmd, checkRestartCmd)
	rootCmd.AddCommand(httpControlCmd)
}
