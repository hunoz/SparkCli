package changepassword

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gtech.dev/spark/cognito"
	"gtech.dev/spark/config"
)

var ChangePasswordCmd = &cobra.Command{
	Use:   "change-password",
	Short: "Change your password. Do not use if needing to reset password",
	Run: func(cmd *cobra.Command, args []string) {
		var configuration *config.CognitoConfig
		config.CheckIfCognitoIsInitialized()
		if config, e := config.GetCognitoConfig(); e != nil {
			fmt.Printf("Error getting config: %s\n", e.Error())
			os.Exit(1)
		} else {
			configuration = config
		}
		cognitoClient := cognito.New(configuration)

		cognitoClient.ChangePassword()
	},
}
