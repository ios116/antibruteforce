package cmd

import (
	"github.com/spf13/cobra"
)

var ResetBucket = &cobra.Command{
	Use:   "bucket",
	Short: "Reset bucket by login or ip",
	Run: func(cmd *cobra.Command, args []string) {
		//container := BuildContainer()

	},
}
