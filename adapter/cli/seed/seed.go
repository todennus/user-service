package seed

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/todennus/shared/middleware"
	"github.com/todennus/user-service/usecase/dto"
	"github.com/todennus/user-service/wiring"
)

var username string
var password string

var Command = &cobra.Command{
	Use:   "seed",
	Short: "Seed the first user",
	Run: func(cmd *cobra.Command, args []string) {
		envPaths, err := cmd.Flags().GetStringArray("env")
		if err != nil {
			panic(err)
		}

		system, err := wiring.InitializeSystem(envPaths...)
		if err != nil {
			panic(err)
		}

		ctx := middleware.WithBasicContext(context.Background(), system.Config)

		resp, err := system.Usecases.UserUsecase.RegisterFirst(ctx, &dto.UserRegisterFirstRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			fmt.Println("Failed:", err)
			return
		}

		fmt.Println("Seed the user successfully")
		fmt.Println("UserID:", resp.User.ID)
	},
}

func init() {
	Command.Flags().StringVarP(&username, "username", "u", "", "username")
	Command.Flags().StringVarP(&password, "password", "p", "", "password")
	Command.MarkFlagRequired("username")
	Command.MarkFlagRequired("password")
}
