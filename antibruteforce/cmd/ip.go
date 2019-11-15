package cmd

import (
	"github.com/spf13/cobra"
)

var CrudIp = &cobra.Command{
	Use:   "ip",
	Short: "crud ip",
	Run: func(cmd *cobra.Command, args []string) {
		//container := BuildContainer()
	},
}
