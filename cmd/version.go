package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of signal-exporter",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nextunit.io Signal Exporter v0.1 -- HEAD")
	},
}
