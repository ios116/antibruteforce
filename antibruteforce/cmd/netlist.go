package cmd

import (
	"antibruteforce/internal/config"
	"antibruteforce/internal/grpcserver"
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Command for adding ip to the blacklist or whitelist",
	//Args:  cobra.MinimumNArgs(1),
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
			var typOfList grpcserver.List
			switch listType {
			case "white":
				typOfList = grpcserver.List_WHITE
			case "black":
				typOfList = grpcserver.List_BLACK
			default:
				fmt.Println("type may by white or black")
				return
			}
			req := &grpcserver.AddIpRequest{Net: ipNet, List: typOfList}
			status, err := server.AddIP(ctx, req)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(status.Ok)
			return
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "The command removes ip from the list",
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
		if err != nil {
			log.Fatal(err)
		}
	},
}

var subnetCmd = &cobra.Command{
	Use:   "subnet",
	Short: "The command look up all subnet",
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
			subn := &grpcserver.GetSubnetRequest{
				Net: ipNet,
			}
			results, err := server.GetSubnet(ctx, subn)
			if err != nil {
				log.Fatal(err)
			}
			for _, net := range results.Nets {
				fmt.Println(net.List, net.Net)
			}
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}
