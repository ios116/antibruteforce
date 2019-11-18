package cmd

import (
	"github.com/spf13/cobra"
)

var resetBucket = &cobra.Command{
	Use:   "bucket",
	//Args: cobra.MinimumNArgs(1),
	Short: "Command resets bucket by login or ip",
	Run: func(cmd *cobra.Command, args []string) {
		//container := BuildContainer()
	},
}
