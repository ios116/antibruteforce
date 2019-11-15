package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var RootCmd = &cobra.Command{
	Use:   "./main",
	Short: "anti brut force",
}

var password string
var login string
var ip string

func init() {
	RootCmd.AddCommand(CrudIp, ResetBucket, )

	// buckets management
	ResetBucket.Flags().StringVarP(&login, "login", "l", "", "bucket login")
	ResetBucket.Flags().StringVarP(&ip, "ip", "i", "", "bucket ip ")

	if err := ResetBucket.MarkFlagRequired("login"); err != nil {
		log.Println(err)
	}

	if err := ResetBucket.MarkFlagRequired("ip"); err != nil {
		log.Println(err)
	}

	// ips management
	CrudIp.Flags().StringVarP(&ip, "ip", "i", "", "ip with mask")

	if err := CrudIp.MarkFlagRequired("ip"); err != nil {
		log.Println(err)
	}
}
