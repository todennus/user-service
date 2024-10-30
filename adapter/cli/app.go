package cli

import (
	"github.com/spf13/cobra"
	"github.com/todennus/user-service/adapter/cli/seed"
)

var Command = &cobra.Command{
	Use:   "cli",
	Short: "The Todennus User CLI",
}

func init() {
	Command.AddCommand(seed.Command)
}
