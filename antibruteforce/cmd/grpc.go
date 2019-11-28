package cmd

import (
	"antibruteforce/internal/grpcserver"
	"log"

	"github.com/spf13/cobra"
)

var grpcRun = &cobra.Command{
	Use:   "grpc",
	Short: "Command to start grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		container := BuildContainer()
		err := container.Invoke(func(serverGRPS *grpcserver.RPCServer) {
			serverGRPS.Start()
		})
		if err != nil {
			log.Println(err)
		}
	},
}
