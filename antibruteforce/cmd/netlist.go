package cmd

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/grpcserver"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
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
		container := BuildContainer()
		err := container.Invoke(func(conf *config.GrpcConf) {
			conn, err := newGrpcConnection(conf)
			defer conn.Close()
			if err != nil {
				log.Fatal(err)
			}
			server := grpcserver.NewAntiBruteForceClient(conn)
			ctx := context.Background()
			req := &grpcserver.DeleteIpRequest{
				Net: ipNet,
			}
			status, err := server.DeleteIP(ctx, req)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(status.Ok)

		})
		log.Fatal(err)
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

var addWhiteList = &cobra.Command{
	Use:   "add",
	Short: "The command adds ip to whitelist",
	//	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add white list", ipNet)
	},
}
