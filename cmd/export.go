package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"nextunit.io/signal-exporter/signal"
)

var (
	exportCmdConfigPath string
	defaultPath         string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export signal chats",
	Long:  "Export given chats into different formats",
	Run: func(cmd *cobra.Command, args []string) {
		var signalObj = signal.New(exportCmdConfigPath)
		defer signalObj.Finish()

		signalObj.Execute()
	},
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		er(err)
	}

	switch runtime.GOOS {
	case "windows":
		defaultPath = home + "/.config/Signal"
	case "darwin":
		defaultPath = home + "/Library/Application Support/Signal"
	case "linux":
		defaultPath = home + "/AppData/Roaming/Signal"
	}

	exportCmd.Flags().StringVarP(&exportCmdConfigPath, "path", "p", defaultPath, "Path to signal")
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
