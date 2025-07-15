package cmd

import (
	"fluidnc-client/internal/config"
	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [file] [destination]",
	Short: "Upload file to FluidNC",
	Long:  "Upload file to FluidNC filesystem or SD card (use --sd flag for SD card).",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		client := fluidnc.NewClient(cfg)
		filePath := args[0]
		destination := ""
		if len(args) > 1 {
			destination = args[1]
		}

		endpoint := "files"
		if sd, _ := cmd.Flags().GetBool("sd"); sd {
			endpoint = "upload"
		}

		return client.UploadFile(filePath, destination, endpoint)
	},
}

func init() {
	uploadCmd.Flags().Bool("sd", false, "Upload to SD card")
	rootCmd.AddCommand(uploadCmd)
}
