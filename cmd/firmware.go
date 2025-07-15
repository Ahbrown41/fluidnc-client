package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"fluidnc-client/internal/config"
	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/cobra"
)

var firmwareCmd = &cobra.Command{
	Use:   "firmware [firmware-file]",
	Short: "Update FluidNC firmware",
	Long:  "Upload and install new firmware via the /updatefw endpoint.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		client := fluidnc.NewClient(cfg)
		firmwarePath := args[0]

		// Confirm before proceeding
		fmt.Printf("WARNING: This will update the FluidNC firmware with file: %s\n", firmwarePath)
		fmt.Print("Are you sure you want to continue? (yes/no): ")

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return fmt.Errorf("failed to read confirmation")
		}

		confirmation := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if confirmation != "yes" && confirmation != "y" {
			fmt.Println("Firmware update cancelled")
			return nil
		}

		if err := client.UpdateFirmware(firmwarePath); err != nil {
			return err
		}

		fmt.Printf("Firmware update initiated successfully\n")
		fmt.Printf("Please wait for the device to restart and complete the update process\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(firmwareCmd)
}
