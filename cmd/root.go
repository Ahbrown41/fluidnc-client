package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "fluidnc-cli",
	Short: "FluidNC Web API Command Line Tool",
	Long: `A command-line tool for FluidNC with real-time monitoring and control.
Supports file uploads, G-code execution, status monitoring, and machine control.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().String("host", "", "FluidNC host address")
	rootCmd.PersistentFlags().Int("port", 0, "FluidNC HTTP port")
	rootCmd.PersistentFlags().Int("websocket-port", 0, "FluidNC WebSocket port")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().String("output", "text", "Output format (text|json)")

	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("websocket_port", rootCmd.PersistentFlags().Lookup("websocket-port"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("output_format", rootCmd.PersistentFlags().Lookup("output"))
}
