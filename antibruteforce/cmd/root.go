package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd root command for DI
var RootCmd = &cobra.Command{
	Use:   "./abf",
	Short: "antibrutforce service",
}

var (
	login, ip, ipNet, listType string
)

func init() {
	RootCmd.AddCommand(addCmd, deleteCmd, bucketCmd, grpcRun, subnetCmd)
	bucketCmd.Flags().StringVarP(&login, "login", "l", "", "bucket login")
	bucketCmd.Flags().StringVarP(&ip, "ip", "i", "", "bucket ip ")
	addCmd.Flags().StringVarP(&listType, "type", "t", "", "type of list may be white/black")
	RootCmd.PersistentFlags().StringVarP(&ipNet, "net", "n", "", "ip with mask example 127.0.0.0/24")
}
