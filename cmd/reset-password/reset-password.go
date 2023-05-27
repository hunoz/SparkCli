package resetpassword

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gtech.dev/spark/cognito"
	"gtech.dev/spark/config"
)

var ResetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset your password",
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckIfCognitoIsInitialized()
		var configuration *config.CognitoConfig
		config.CheckIfCognitoIsInitialized()
		if config, e := config.GetCognitoConfig(); e != nil {
			fmt.Printf("Error getting config: %s\n", e.Error())
			os.Exit(1)
		} else {
			configuration = config
		}
		cognitoClient := cognito.New(configuration)

		cognitoClient.ResetPassword()
	},
}
