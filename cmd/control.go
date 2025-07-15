package cmd

import (
	"fluidnc-client/internal/config"
	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/cobra"
)

var controlCmd = &cobra.Command{
	Use:   "control",
	Short: "Machine control commands",
}

var holdCmd = &cobra.Command{
	Use:   "hold",
	Short: "Send feed hold",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		client := fluidnc.NewClient(cfg)
		if err := client.Connect(); err != nil {
			return err
		}
		defer client.Disconnect()
		return client.FeedHold()
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Send cycle start",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		client := fluidnc.NewClient(cfg)
		if err := client.Connect(); err != nil {
			return err
		}
		defer client.Disconnect()
		return client.CycleStart()
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Send soft reset",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		client := fluidnc.NewClient(cfg)
		if err := client.Connect(); err != nil {
			return err
		}
		defer client.Disconnect()
		return client.SoftReset()
	},
}

var homeCmd = &cobra.Command{
	Use:   "home",
	Short: "Home machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		client := fluidnc.NewClient(cfg)
		if err := client.Connect(); err != nil {
			return err
		}
		defer client.Disconnect()
		return client.Home()
	},
}

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.LoadConfig()
		client := fluidnc.NewClient(cfg)
		if err := client.Connect(); err != nil {
			return err
		}
		defer client.Disconnect()
		return client.Unlock()
	},
}

func init() {
	controlCmd.AddCommand(holdCmd, startCmd, resetCmd, homeCmd, unlockCmd)
	rootCmd.AddCommand(controlCmd)
}
