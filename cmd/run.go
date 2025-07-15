package cmd

import (
	"fluidnc-client/internal/config"
	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Run G-code file with monitoring",
	Long:  "Execute G-code file line by line with optional real-time status monitoring.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		client := fluidnc.NewClient(cfg)
		monitor, _ := cmd.Flags().GetBool("monitor")

		return client.RunGCodeFile(args[0], monitor)
	},
}

func init() {
	runCmd.Flags().Bool("monitor", true, "Enable real-time status monitoring")
	rootCmd.AddCommand(runCmd)
}
