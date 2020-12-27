package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "signal-exporter",
	Short: "Signal exporter exports your signal chats.",
	Long:  "Signal exporter exports your signal chats into different output formats",
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(exportCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
