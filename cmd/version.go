package cmd

import (
	"fmt"

	"github.com/ediblesushi/helios/pkg/printing"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(printing.VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
