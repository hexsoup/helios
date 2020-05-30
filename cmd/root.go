package cmd

import (
	"fmt"
	"os"

	"github.com/ediblesushi/helios/pkg/printing"
	"github.com/ediblesushi/helios/pkg/scanning"
	"github.com/spf13/cobra"
)

var optionsStr = make(map[string]string)
var optionsBool = make(map[string]bool)
var target, ports string
var verbose bool

var rootCmd = &cobra.Command{
	Use: "helios",
	Run: func(cmd *cobra.Command, args []string) {
		scanning.Scan(optionsStr, optionsBool)
	},
}

// Execute will be called to execute the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Str
	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "target to be scanned (required)")
	rootCmd.MarkPersistentFlagRequired("target")
	rootCmd.PersistentFlags().StringVarP(&ports, "ports", "p", "", "ports to be scanned")

	//Int

	//Bool
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output")

	printing.Banner()
}

func initConfig() {
	optionsStr["target"] = target
	optionsStr["ports"] = ports
	optionsBool["verbose"] = verbose
}
