package cmd

import (
	"fmt"

	"fluidnc-client/internal/config"
	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/cobra"
)

var httpCmdCmd = &cobra.Command{
	Use:   "http-cmd [command]",
	Short: "Send command via HTTP",
	Long:  "Send a command to FluidNC via HTTP POST to /command endpoint.  For example, use 'http-cmd \\$S' to send the '$S' command.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		client := fluidnc.NewClient(cfg)
		silent, _ := cmd.Flags().GetBool("silent")

		response, err := client.SendHTTPCommand(args[0], silent)
		if err != nil {
			return err
		}

		if !silent || cfg.Verbose {
			fmt.Println(response)
		}
		return nil
	},
}

func init() {
	httpCmdCmd.Flags().Bool("silent", false, "Send command silently (use /command_silent endpoint)")
	rootCmd.AddCommand(httpCmdCmd)
}
