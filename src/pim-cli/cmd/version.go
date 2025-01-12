package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version number of pim cli",
	Long:    `Print the version number of pim cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1")
	},
}
