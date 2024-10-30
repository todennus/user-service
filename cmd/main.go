package main

import (
	"github.com/spf13/cobra"
	"github.com/todennus/user-service/adapter/cli"
	"github.com/todennus/user-service/cmd/grpc"
	"github.com/todennus/user-service/cmd/rest"
)

var rootCommand = &cobra.Command{
	Use:   "todennus",
	Short: "todennus is an Identity, OpenID Connect, and OAuth2 provider",
}

func main() {
	rootCommand.PersistentFlags().StringArray("env", []string{".env"}, "environment file paths")
	rootCommand.AddCommand(rest.Command)
	rootCommand.AddCommand(grpc.Command)
	rootCommand.AddCommand(cli.Command)

	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}
