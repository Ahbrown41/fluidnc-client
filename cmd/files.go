package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"fluidnc-client/internal/config"
	"fluidnc-client/internal/fluidnc"
	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "File operations",
	Long:  "List, upload, and manage files on FluidNC filesystem and SD card.",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files on FluidNC filesystem",
	Long:  "List all files and directories on the FluidNC local filesystem.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}

		client := fluidnc.NewClient(cfg)
		fileList, err := client.ListFiles()
		if err != nil {
			return err
		}

		if cfg.OutputFormat == "json" {
			jsonOutput, _ := json.MarshalIndent(fileList, "", "  ")
			fmt.Println(string(jsonOutput))
		} else {
			fmt.Printf("Files in %s:\n", fileList.Path)
			fmt.Printf("%-30s %-10s %s\n", "Name", "Type", "Size")
			fmt.Println(strings.Repeat("-", 50))
			for _, file := range fileList.Files {
				sizeStr := ""
				if file.Type == "file" && file.Size > 0 {
					sizeStr = fmt.Sprintf("%d bytes", file.Size)
				}
				fmt.Printf("%-30s %-10s %s\n", file.Name, file.Type, sizeStr)
			}
		}
		return nil
	},
}

var uploadLocalCmd = &cobra.Command{
	Use:   "upload-local [file] [destination]",
	Short: "Upload file to local filesystem",
	Long:  "Upload a file to the FluidNC local filesystem via the /files endpoint.",
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

		if err := client.UploadToLocalFS(filePath, destination); err != nil {
			return err
		}

		fmt.Printf("File uploaded successfully to local filesystem\n")
		return nil
	},
}

var uploadSDCmd = &cobra.Command{
	Use:   "upload-sd [file] [destination]",
	Short: "Upload file to SD card",
	Long:  "Upload a file directly to the SD card via the /upload endpoint.",
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

		if err := client.UploadToSD(filePath, destination); err != nil {
			return err
		}

		fmt.Printf("File uploaded successfully to SD card\n")
		return nil
	},
}

func init() {
	filesCmd.AddCommand(listCmd, uploadLocalCmd, uploadSDCmd)
	rootCmd.AddCommand(filesCmd)
}
