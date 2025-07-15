package cmd

import (
	"fluidnc-client/internal/config"
	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactive session",
	Long:  "Start an interactive WebSocket session with FluidNC.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		client := fluidnc.NewClient(cfg)
		return client.InteractiveMode()
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
