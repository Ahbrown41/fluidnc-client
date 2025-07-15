package cmd

import (
	"context"

	"fluidnc-client/internal/config"
	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Monitor FluidNC status",
	Long:  "Continuously monitor FluidNC status with real-time updates.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		client := fluidnc.NewClient(cfg)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		return client.MonitorStatus(ctx, func(status *fluidnc.FluidNCStatus) {
			client.DisplayStatus(status)
		})
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
