package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd root command for DI
var RootCmd = &cobra.Command{
	Use:   "abf",
	Short: "anti brut force",
}

var (
	login, ip, ipNet string
)

func init() {

	RootCmd.AddCommand(blackList, whiteList, resetBucket, grpcRun)
	// buckets management
	resetBucket.Flags().StringVarP(&login, "login", "l", "", "bucket login")
	resetBucket.Flags().StringVarP(&ip, "ip", "i", "", "bucket ip ")
	whiteList.AddCommand(addWhiteList, deleteWhiteList)
	blackList.AddCommand(addBlackList, deleteBlackList)
	RootCmd.PersistentFlags().StringVarP(&ipNet, "net", "n", "", "action with net list maybe delete or add")

}
