package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var blackList = &cobra.Command{
	Use:   "blacklist",
	Short: "Command choice the blacklist",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("blacklist= ", ipNet)
	},
}

var deleteBlackList = &cobra.Command{
	Use:   "delete",
	Short: "The command removes ip from the blacklist",
	//	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete black list", ipNet)
	},
}

var addBlackList = &cobra.Command{
	Use:   "add",
	Short: "The command adds ip to the blacklist",
	//	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add black list", ipNet)
	},
}

var whiteList = &cobra.Command{
	Use:   "whitelist",
	Short: "Command choice the whitelist",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("white=", ipNet)
	},
}

var deleteWhiteList = &cobra.Command{
	Use:   "delete",
	Short: "The command removes ip from the whitelist",
	//	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete white list", ipNet)
	},
}

var addWhiteList = &cobra.Command{
	Use:   "add",
	Short: "The command adds ip to whitelist",
	//	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add white list", ipNet)
	},
}
